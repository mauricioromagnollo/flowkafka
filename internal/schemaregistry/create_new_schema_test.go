package schemaregistry

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// registryStub is a fake Confluent Schema Registry that serves the routes used by
// CreateNewSchema and records every schema registration request it receives.
type registryStub struct {
	mutex sync.Mutex

	// latestStatus and latestBody configure the GET /subjects/{subject}/versions/latest response.
	latestStatus int
	latestBody   []byte

	// postStatus configures the POST /subjects/{subject}/versions response.
	postStatus int

	postCount  int
	postBodies []map[string]any
}

func (s *registryStub) handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/versions/latest"):
		s.mutex.Lock()
		status, body := s.latestStatus, s.latestBody
		s.mutex.Unlock()

		w.WriteHeader(status)
		_, err := w.Write(body)
		Expect(err).ToNot(HaveOccurred())

	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/versions"):
		raw, err := io.ReadAll(r.Body)
		Expect(err).ToNot(HaveOccurred())

		var parsed map[string]any
		Expect(json.Unmarshal(raw, &parsed)).To(Succeed())

		s.mutex.Lock()
		s.postCount++
		s.postBodies = append(s.postBodies, parsed)
		status := s.postStatus
		s.mutex.Unlock()

		w.WriteHeader(status)
		if status == http.StatusOK {
			_, err = w.Write([]byte(`{"id":42}`))
		} else {
			_, err = w.Write([]byte(`{"error_code":42201,"message":"any error reason returned by the registry"}`))
		}
		Expect(err).ToNot(HaveOccurred())

	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/schemas/ids/"):
		// CreateSchema fetches the newly registered schema by id right after the POST.
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(newSchemaResponseBody("any-subject-value", testAvroSchema, 1, 42))
		Expect(err).ToNot(HaveOccurred())

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *registryStub) totalPosts() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.postCount
}

func (s *registryStub) lastPostBody() map[string]any {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	Expect(s.postBodies).ToNot(BeEmpty())

	return s.postBodies[len(s.postBodies)-1]
}

var _ = Describe("CreateNewSchema", func() {
	var stub *registryStub
	var server *httptest.Server
	var client *Client

	// notFoundBody is the payload returned by the registry for an unknown subject.
	notFoundBody := []byte(`{"error_code":40401,"message":"Subject 'any-subject-value' not found."}`)

	BeforeEach(func() {
		stub = &registryStub{
			latestStatus: http.StatusOK,
			latestBody:   newSchemaResponseBody("any-subject-value", testAvroSchema, 1, 42),
			postStatus:   http.StatusOK,
		}
		server = httptest.NewServer(http.HandlerFunc(stub.handler))

		// A fresh client per spec avoids contamination from the srclient internal cache.
		client = NewSchemaRegistry(Config{Endpoint: server.URL})
	})

	AfterEach(func() {
		server.Close()
	})

	It("should return nil when the latest schema is already identical", func() {
		err := client.CreateNewSchema("any-subject-value", testAvroSchema, SchemaTypeAvro)

		Expect(err).ToNot(HaveOccurred())
		Expect(stub.totalPosts()).To(Equal(0))
	})

	It("should return nil when the existing schema differs and a new version is created", func() {
		stub.latestBody = newSchemaResponseBody("any-subject-value", testOtherAvroSchema, 1, 42)

		err := client.CreateNewSchema("any-subject-value", testAvroSchema, SchemaTypeAvro)

		Expect(err).ToNot(HaveOccurred())
		Expect(stub.totalPosts()).To(Equal(1))
		Expect(stub.lastPostBody()).To(HaveKeyWithValue("schema", testAvroSchema))
	})

	It("should return nil when the subject does not exist yet", func() {
		stub.latestStatus = http.StatusNotFound
		stub.latestBody = notFoundBody

		err := client.CreateNewSchema("any-subject-value", testAvroSchema, SchemaTypeAvro)

		Expect(err).ToNot(HaveOccurred())
		Expect(stub.totalPosts()).To(Equal(1))
	})

	It("should return error when the schema registration fails", func() {
		stub.latestStatus = http.StatusNotFound
		stub.latestBody = notFoundBody
		stub.postStatus = http.StatusUnprocessableEntity

		err := client.CreateNewSchema("any-subject-value", testAvroSchema, SchemaTypeAvro)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("any error reason returned by the registry"))
		Expect(stub.totalPosts()).To(Equal(1))
	})

	It("should return error when the endpoint is unreachable", func() {
		unreachableClient := NewSchemaRegistry(Config{Endpoint: "http://127.0.0.1:1"})

		err := unreachableClient.CreateNewSchema("any-subject-value", testAvroSchema, SchemaTypeAvro)

		Expect(err).To(HaveOccurred())
	})

	It("should return error when the schema type is not supported", func() {
		stub.latestStatus = http.StatusNotFound
		stub.latestBody = notFoundBody

		err := client.CreateNewSchema("any-subject-value", testAvroSchema, SchemaType("any-invalid-type"))

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid schema type"))
		Expect(stub.totalPosts()).To(Equal(0))
	})

	// srclient omits the schemaType field for Avro, since it is the registry default.
	It("should return nil when SchemaTypeAvro is propagated to the registry", func() {
		stub.latestStatus = http.StatusNotFound
		stub.latestBody = notFoundBody

		err := client.CreateNewSchema("any-subject-value", testAvroSchema, SchemaTypeAvro)

		Expect(err).ToNot(HaveOccurred())
		Expect(stub.lastPostBody()).ToNot(HaveKey("schemaType"))
	})

	It("should return nil when SchemaTypeJSON is propagated to the registry", func() {
		stub.latestStatus = http.StatusNotFound
		stub.latestBody = notFoundBody

		err := client.CreateNewSchema("any-subject-value", testAvroSchema, SchemaTypeJSON)

		Expect(err).ToNot(HaveOccurred())
		Expect(stub.lastPostBody()).To(HaveKeyWithValue("schemaType", "JSON"))
	})

	It("should return nil when SchemaTypeProtobuf is propagated to the registry", func() {
		stub.latestStatus = http.StatusNotFound
		stub.latestBody = notFoundBody

		err := client.CreateNewSchema("any-subject-value", testAvroSchema, SchemaTypeProtobuf)

		Expect(err).ToNot(HaveOccurred())
		Expect(stub.lastPostBody()).To(HaveKeyWithValue("schemaType", "PROTOBUF"))
	})
})
