package methods1

type Counter struct {
	n int
}

func (c *Counter) Inc() {
	c.n++
}

func (c Counter) Value() int {
	return c.n
}
