package producer

// HasSchemaRegistry returns true if the client has a schema registry configured.
func (c *producerClient) HasSchemaRegistry() bool {
	return c.schemaRegistry != nil
}
