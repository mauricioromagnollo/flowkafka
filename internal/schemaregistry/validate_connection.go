package schemaregistry

import (
	"fmt"
)

// ValidateConnection validates the connection to the Schema Registry
// by attempting to retrieve the list of subjects.
func (c *Client) ValidateConnection() error {
	_, err := c.client.GetSubjects()
	if err != nil {
		return fmt.Errorf("failed to connect to schema registry: %w", err)
	}

	return nil
}
