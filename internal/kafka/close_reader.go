package kafka

// CloseReader closes the reader for the given topic, if it exists.
func (c *Client) CloseReader(topic string) error {
	c.mu.Lock()
	r, ok := c.readers[topic]
	if ok {
		delete(c.readers, topic)
	}
	c.mu.Unlock()

	if !ok {
		return nil
	}
	return r.Close()
}
