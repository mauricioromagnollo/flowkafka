package schemaregistry

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSchemaRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SchemaRegistry Suite")
}

// testAvroSchema is a valid Avro schema used as fixture across the specs.
const testAvroSchema = `{"type":"record","name":"Test","fields":[{"name":"id","type":"string"}]}`

// testOtherAvroSchema is a valid Avro schema that differs from testAvroSchema.
const testOtherAvroSchema = `{"type":"record","name":"Test","fields":[{"name":"name","type":"string"}]}`

// newSchemaResponseBody builds the JSON body returned by the Confluent Schema Registry
// for schema lookups, properly escaping the schema string inside the JSON payload.
func newSchemaResponseBody(subject, schema string, version, id int) []byte {
	body, err := json.Marshal(map[string]any{
		"subject": subject,
		"version": version,
		"id":      id,
		"schema":  schema,
	})
	Expect(err).ToNot(HaveOccurred())

	return body
}

// basicAuthValue builds the base64 encoded value of an HTTP Basic Authorization header.
func basicAuthValue(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}
