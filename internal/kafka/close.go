package kafka

// Close closes the Kafka client, including all readers and the writer.
func (c *Client) Close() error {
	c.mu.Lock()
	readers := make([]func() error, 0, len(c.readers))
	for topic, r := range c.readers {
		rr := r
		readers = append(readers, func() error { return rr.Close() })
		delete(c.readers, topic)
	}
	c.mu.Unlock()

	for _, closeFn := range readers {
		_ = closeFn()
	}

	if c.writer != nil {
		return c.writer.Close()
	}
	return nil
}
