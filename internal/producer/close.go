package producer

func (c *producerClient) Close() error {
	if c.writer != nil {
		return c.writer.Close()
	}

	return nil
}
