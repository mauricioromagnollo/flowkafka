package schemaregistry

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidateConnection", func() {
	It("should return nil when the schema registry returns the list of subjects", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.Method).To(Equal(http.MethodGet))
			Expect(r.URL.Path).To(Equal("/subjects"))

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`["any-subject-value"]`))
			Expect(err).ToNot(HaveOccurred())
		}))
		defer server.Close()

		client := NewSchemaRegistry(Config{Endpoint: server.URL})

		err := client.ValidateConnection()

		Expect(err).ToNot(HaveOccurred())
	})

	It("should return error when the schema registry returns an internal server error", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`{"error_code":50001,"message":"any error reason returned by the registry"}`))
			Expect(err).ToNot(HaveOccurred())
		}))
		defer server.Close()

		client := NewSchemaRegistry(Config{Endpoint: server.URL})

		err := client.ValidateConnection()

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to connect to schema registry"))
	})

	It("should return error when the schema registry returns an invalid payload", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`not-a-valid-json`))
			Expect(err).ToNot(HaveOccurred())
		}))
		defer server.Close()

		client := NewSchemaRegistry(Config{Endpoint: server.URL})

		err := client.ValidateConnection()

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to connect to schema registry"))
	})

	It("should return error when the endpoint is unreachable", func() {
		client := NewSchemaRegistry(Config{Endpoint: "http://127.0.0.1:1"})

		err := client.ValidateConnection()

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to connect to schema registry"))
	})
})
