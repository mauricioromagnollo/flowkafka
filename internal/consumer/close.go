package consumer

// Close closes the consumer and releases any resources it holds.
func (c *consumerClient) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}

	return nil
}
