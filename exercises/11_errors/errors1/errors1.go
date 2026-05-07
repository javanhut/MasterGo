package errors1

// I AM NOT DONE

// Divide should return (a/b, nil) when b != 0,
// and (0, error) when b == 0. The error message must contain the
// substring "divide by zero".

func Divide(a, b int) (int, error) {
	return a / b, nil
}
