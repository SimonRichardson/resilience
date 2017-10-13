package breaker_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/SimonRichardson/resilience/breaker"
)

func TestBreaker(t *testing.T) {
	t.Parallel()

	t.Run("func called", func(t *testing.T) {
		called := false

		breaker := breaker.New(1, time.Second)
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
	})

	t.Run("func returns error", func(t *testing.T) {
		called := false

		breaker := breaker.New(1, time.Second)
		err := breaker.Run(func() error {
			called = true
			return errors.New("bad")
		})

		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, err != nil; expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})

	t.Run("second call isn't successful", func(t *testing.T) {
		called := false

		breaker := breaker.New(1, time.Second)
		breaker.Run(func() error {
			return errors.New("bad")
		})

		err := breaker.Run(func() error {
			called = true
			return nil
		})

		if expected, actual := false, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, err != nil; expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})

	t.Run("second call is successful after timeout", func(t *testing.T) {
		called := false

		expiryTimeout := 100 * time.Millisecond

		breaker := breaker.New(1, expiryTimeout)
		breaker.Run(func() error {
			return errors.New("bad")
		})

		time.Sleep(expiryTimeout)

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
	})
}

func Example() {
	circuit := breaker.New(10, 10*time.Millisecond)
	err := circuit.Run(func() error {
		// communicate with an external service
		return nil
	})

	switch {
	case err == nil:
		fmt.Println("success!")
	case breaker.ErrBreakerOpen(err):
		fmt.Println("breaker is open")
	default:
		fmt.Println("other error")
	}

	// Output:
	// success!
}
