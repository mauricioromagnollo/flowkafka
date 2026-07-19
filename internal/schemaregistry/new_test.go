package schemaregistry

import (
	"net/http"
	"net/http/httptest"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewSchemaRegistry", func() {
	// authRecorder starts a fake schema registry that records the Authorization
	// header of the last received request.
	newAuthRecorderServer := func() (*httptest.Server, func() string) {
		var mutex sync.Mutex
		var lastAuthorization string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mutex.Lock()
			lastAuthorization = r.Header.Get("Authorization")
			mutex.Unlock()

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`[]`))
			Expect(err).ToNot(HaveOccurred())
		}))

		return server, func() string {
			mutex.Lock()
			defer mutex.Unlock()

			return lastAuthorization
		}
	}

	It("should return a client when a valid config is provided", func() {
		config := Config{Endpoint: "http://localhost:8081"}

		client := NewSchemaRegistry(config)

		Expect(client).ToNot(BeNil())
		Expect(client.config).To(Equal(config))
		Expect(client.client).ToNot(BeNil())
		Expect(client.client.GetSchemaRegistryURL()).To(Equal("http://localhost:8081"))
	})

	It("should return nil authorization header when username is empty", func() {
		server, lastAuthorization := newAuthRecorderServer()
		defer server.Close()

		client := NewSchemaRegistry(Config{
			Endpoint:     server.URL,
			SaslUsername: "",
			SaslPassword: "any-password",
		})

		err := client.ValidateConnection()

		Expect(err).ToNot(HaveOccurred())
		Expect(lastAuthorization()).To(BeEmpty())
	})

	It("should return nil authorization header when password is empty", func() {
		server, lastAuthorization := newAuthRecorderServer()
		defer server.Close()

		client := NewSchemaRegistry(Config{
			Endpoint:     server.URL,
			SaslUsername: "any-username",
			SaslPassword: "",
		})

		err := client.ValidateConnection()

		Expect(err).ToNot(HaveOccurred())
		Expect(lastAuthorization()).To(BeEmpty())
	})

	It("should return basic auth header when username and password are provided", func() {
		server, lastAuthorization := newAuthRecorderServer()
		defer server.Close()

		client := NewSchemaRegistry(Config{
			Endpoint:     server.URL,
			SaslUsername: "any-username",
			SaslPassword: "any-password",
		})

		err := client.ValidateConnection()

		Expect(err).ToNot(HaveOccurred())

		expected := "Basic " + basicAuthValue("any-username", "any-password")
		Expect(lastAuthorization()).To(Equal(expected))
	})
})
