package schemaregistry

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetLatestSchema", func() {
	It("should return the latest schema when the subject exists", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.Method).To(Equal(http.MethodGet))
			Expect(r.URL.Path).To(Equal("/subjects/any-subject-value/versions/latest"))

			w.WriteHeader(http.StatusOK)
			_, err := w.Write(newSchemaResponseBody("any-subject-value", testAvroSchema, 3, 42))
			Expect(err).ToNot(HaveOccurred())
		}))
		defer server.Close()

		client := NewSchemaRegistry(Config{Endpoint: server.URL})

		schema, err := client.GetLatestSchema("any-subject-value")

		Expect(err).ToNot(HaveOccurred())
		Expect(schema).ToNot(BeNil())
		Expect(schema.ID()).To(Equal(42))
		Expect(schema.Version()).To(Equal(3))
		Expect(schema.Schema()).To(Equal(testAvroSchema))
	})

	It("should return error when the subject does not exist", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"error_code":40401,"message":"Subject 'any-subject-value' not found."}`))
			Expect(err).ToNot(HaveOccurred())
		}))
		defer server.Close()

		client := NewSchemaRegistry(Config{Endpoint: server.URL})

		schema, err := client.GetLatestSchema("any-subject-value")

		Expect(err).To(HaveOccurred())
		Expect(schema).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("should return error when the schema registry returns an invalid payload", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`not-a-valid-json`))
			Expect(err).ToNot(HaveOccurred())
		}))
		defer server.Close()

		client := NewSchemaRegistry(Config{Endpoint: server.URL})

		schema, err := client.GetLatestSchema("any-subject-value")

		Expect(err).To(HaveOccurred())
		Expect(schema).To(BeNil())
	})

	It("should return error when the endpoint is unreachable", func() {
		client := NewSchemaRegistry(Config{Endpoint: "http://127.0.0.1:1"})

		schema, err := client.GetLatestSchema("any-subject-value")

		Expect(err).To(HaveOccurred())
		Expect(schema).To(BeNil())
	})
})
