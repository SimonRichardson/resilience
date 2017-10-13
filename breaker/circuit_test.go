package breaker_test

import (
	"testing"
	"time"

	"github.com/SimonRichardson/resilience/breaker"
)

func TestBreaker(t *testing.T) {

	called := false

	breaker := breaker.New(1, 1, time.Second)
	err := breaker.Run(func() error {
		called = true
		return nil
	})

	if expected, actual := true, called; expected != actual {
		t.Errorf("expected: %t, actual: %t", expected, actual)
	}
	if expected, actual := true, err == nil; expected != actual {
		t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
	}
}
