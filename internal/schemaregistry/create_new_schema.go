package schemaregistry

import "github.com/riferrei/srclient"

// CreateNewSchema creates a new schema for a given subject in the schema registry.
func (c *Client) CreateNewSchema(subject, schema string, schemaType SchemaType) error {
	// Check if schema already exists
	existingSchema, err := c.client.GetLatestSchema(subject)
	if err == nil && existingSchema != nil {
		if existingSchema.Schema() == schema {
			return nil
		}
	}

	// Schema doesn't exist or differs, create new version
	_, err = c.client.CreateSchema(subject, schema, srclient.SchemaType(schemaType))
	if err != nil {
		return err
	}
	return nil
}
