package errors1

import (
	"strings"
	"testing"
)

func TestDivideOK(t *testing.T) {
	q, err := Divide(10, 2)
	if err != nil || q != 5 {
		t.Fatalf("Divide(10,2)=(%d,%v); want (5,nil)", q, err)
	}
}

func TestDivideZero(t *testing.T) {
	_, err := Divide(10, 0)
	if err == nil {
		t.Fatal("Divide(10,0) returned nil error")
	}
	if !strings.Contains(err.Error(), "divide by zero") {
		t.Errorf("error %q does not mention 'divide by zero'", err)
	}
}
