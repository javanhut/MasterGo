package methods1

// I AM NOT DONE

// Counter must support Inc() to increment and Value() to read.
// Make Inc actually mutate the receiver — the test calls it through a *Counter.

type Counter struct {
	n int
}

func (c Counter) Inc() {
	c.n++
}

func (c Counter) Value() int {
	return c.n
}
