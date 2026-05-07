package methods1

import "testing"

func TestCounter(t *testing.T) {
	c := &Counter{}
	c.Inc()
	c.Inc()
	c.Inc()
	if c.Value() != 3 {
		t.Errorf("Value=%d; want 3", c.Value())
	}
}
