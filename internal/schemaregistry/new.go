package schemaregistry

import "github.com/riferrei/srclient"

// Client is the internal Schema Registry client that wraps the srclient.SchemaRegistryClient.
type Client struct {
	config Config
	client *srclient.SchemaRegistryClient
}

// NewSchemaRegistry creates and configures a new SchemaRegistry client instance.
// It initializes the client with the provided endpoint and optionally sets up
// SASL authentication credentials if both username and password are provided in the config.
func NewSchemaRegistry(config Config) *Client {
	sr := srclient.CreateSchemaRegistryClient(config.Endpoint)

	if config.SaslUsername != "" && config.SaslPassword != "" {
		sr.SetCredentials(config.SaslUsername, config.SaslPassword)
	}

	return &Client{
		config: config,
		client: sr,
	}
}
