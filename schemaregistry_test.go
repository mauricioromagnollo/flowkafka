package flowkafka

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TESTABILITY NOTE: schemaRegistryWrapper holds a concrete *schemaregistry.Client
// instead of the schemaregistry.SchemaRegistry interface, so there is no seam to
// inject a stub. Delegation is therefore proven end to end against a fake HTTP
// schema registry: each spec asserts the request the wrapper caused the internal
// client to issue, plus the value it returned unchanged.
var _ = Describe("SchemaRegistry", func() {
	const anyAvroSchema = `{"type":"record","name":"Test","fields":[{"name":"id","type":"string"}]}`

	type recordedRequest struct {
		Method string
		Path   string
	}

	// newFakeRegistry starts a fake schema registry that records every received
	// request and answers with the handler provided by the spec.
	newFakeRegistry := func(handler http.HandlerFunc) (*httptest.Server, func() []recordedRequest) {
		var mutex sync.Mutex
		var requests []recordedRequest

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mutex.Lock()
			requests = append(requests, recordedRequest{Method: r.Method, Path: r.URL.Path})
			mutex.Unlock()

			handler(w, r)
		}))

		return server, func() []recordedRequest {
			mutex.Lock()
			defer mutex.Unlock()

			return append([]recordedRequest(nil), requests...)
		}
	}

	writeJSON := func(w http.ResponseWriter, status int, body any) {
		payload, err := json.Marshal(body)
		Expect(err).ToNot(HaveOccurred())

		w.WriteHeader(status)
		_, err = w.Write(payload)
		Expect(err).ToNot(HaveOccurred())
	}

	Describe("NewSchemaRegistry", func() {
		It("should return a schema registry when a config is provided", func() {
			registry := NewSchemaRegistry(SchemaRegistryConfig{
				Endpoint:     "http://localhost:8081",
				SaslUsername: "any-username",
				SaslPassword: "any-password",
			})

			Expect(registry).ToNot(BeNil())

			wrapper, ok := registry.(*schemaRegistryWrapper)
			Expect(ok).To(BeTrue())
			Expect(wrapper.client).ToNot(BeNil())
		})
	})

	Describe("ValidateConnection", func() {
		It("should return nil when the registry returns the list of subjects", func() {
			server, requests := newFakeRegistry(func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(w, http.StatusOK, []string{"any-subject-value"})
			})
			defer server.Close()

			registry := NewSchemaRegistry(SchemaRegistryConfig{Endpoint: server.URL})

			err := registry.ValidateConnection()

			Expect(err).ToNot(HaveOccurred())
			Expect(requests()).To(ContainElement(recordedRequest{Method: http.MethodGet, Path: "/subjects"}))
		})

		It("should return error when the registry is unreachable", func() {
			registry := NewSchemaRegistry(SchemaRegistryConfig{Endpoint: "http://127.0.0.1:1"})

			err := registry.ValidateConnection()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to connect to schema registry"))
		})
	})

	Describe("GetLatestSchema", func() {
		It("should return the schema when the registry returns it", func() {
			server, requests := newFakeRegistry(func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(w, http.StatusOK, map[string]any{
					"subject": "any-subject-value",
					"version": 1,
					"id":      4242,
					"schema":  anyAvroSchema,
				})
			})
			defer server.Close()

			registry := NewSchemaRegistry(SchemaRegistryConfig{Endpoint: server.URL})

			schema, err := registry.GetLatestSchema("any-subject-value")

			Expect(err).ToNot(HaveOccurred())
			Expect(schema).ToNot(BeNil())
			Expect(schema.ID()).To(Equal(4242))
			Expect(schema.Schema()).To(Equal(anyAvroSchema))
			Expect(requests()).To(ContainElement(recordedRequest{
				Method: http.MethodGet,
				Path:   "/subjects/any-subject-value/versions/latest",
			}))
		})

		It("should return error when the registry returns an error", func() {
			server, _ := newFakeRegistry(func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(w, http.StatusNotFound, map[string]any{
					"error_code": 40401,
					"message":    "any error reason returned by the registry",
				})
			})
			defer server.Close()

			registry := NewSchemaRegistry(SchemaRegistryConfig{Endpoint: server.URL})

			schema, err := registry.GetLatestSchema("any-subject-value")

			Expect(err).To(HaveOccurred())
			Expect(schema).To(BeNil())
		})
	})

	Describe("CreateNewSchema", func() {
		It("should return nil when the schema is already registered with the same content", func() {
			server, requests := newFakeRegistry(func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(w, http.StatusOK, map[string]any{
					"subject": "any-subject-value",
					"version": 1,
					"id":      4242,
					"schema":  anyAvroSchema,
				})
			})
			defer server.Close()

			registry := NewSchemaRegistry(SchemaRegistryConfig{Endpoint: server.URL})

			err := registry.CreateNewSchema("any-subject-value", anyAvroSchema, SchemaType("AVRO"))

			Expect(err).ToNot(HaveOccurred())
			Expect(requests()).ToNot(ContainElement(recordedRequest{
				Method: http.MethodPost,
				Path:   "/subjects/any-subject-value/versions",
			}))
		})

		It("should return nil when the schema is created", func() {
			server, requests := newFakeRegistry(func(w http.ResponseWriter, r *http.Request) {
				switch {
				// srclient posts the new schema and then re-fetches it by ID.
				case r.Method == http.MethodPost:
					writeJSON(w, http.StatusOK, map[string]any{"id": 4243})
				case strings.HasPrefix(r.URL.Path, "/schemas/ids/"):
					writeJSON(w, http.StatusOK, map[string]any{
						"id":     4243,
						"schema": anyAvroSchema,
					})
				default:
					writeJSON(w, http.StatusNotFound, map[string]any{
						"error_code": 40401,
						"message":    "any error reason returned by the registry",
					})
				}
			})
			defer server.Close()

			registry := NewSchemaRegistry(SchemaRegistryConfig{Endpoint: server.URL})

			err := registry.CreateNewSchema("any-subject-value", anyAvroSchema, SchemaType("AVRO"))

			Expect(err).ToNot(HaveOccurred())
			Expect(requests()).To(ContainElement(recordedRequest{
				Method: http.MethodPost,
				Path:   "/subjects/any-subject-value/versions",
			}))
		})

		It("should return error when the registry fails to create the schema", func() {
			server, _ := newFakeRegistry(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					writeJSON(w, http.StatusConflict, map[string]any{
						"error_code": 409,
						"message":    "any error reason returned by the registry",
					})

					return
				}

				writeJSON(w, http.StatusNotFound, map[string]any{
					"error_code": 40401,
					"message":    "any error reason returned by the registry",
				})
			})
			defer server.Close()

			registry := NewSchemaRegistry(SchemaRegistryConfig{Endpoint: server.URL})

			err := registry.CreateNewSchema("any-subject-value", anyAvroSchema, SchemaType("AVRO"))

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("any error reason returned by the registry"))
		})
	})
})
