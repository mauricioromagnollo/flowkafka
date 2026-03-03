package schemaregistry

import "github.com/riferrei/srclient"

// GetLatestSchema retrieves the latest schema for a given subject from the schema registry.
func (c *Client) GetLatestSchema(subject string) (*srclient.Schema, error) {
	schema, err := c.client.GetLatestSchema(subject)
	if err != nil {
		return nil, err
	}
	return schema, nil
}
