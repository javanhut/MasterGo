package channels1

// I AM NOT DONE

// Produce should send the integers 1..n on the returned channel,
// then close it so the caller can `for v := range ch { ... }`.
// Do the sending in a goroutine so the function can return immediately.

func Produce(n int) <-chan int {
	ch := make(chan int)
	close(ch)
	return ch
}
