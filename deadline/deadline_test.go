package deadline_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/SimonRichardson/resilience/deadline"
)

func TestDeadline(t *testing.T) {
	t.Parallel()

	t.Run("func called", func(t *testing.T) {
		called := false

		timeout := deadline.New(time.Second)
		err := timeout.Run(func(cancel func()) error {
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

		timeout := deadline.New(time.Second)
		err := timeout.Run(func(cancel func()) error {
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

	t.Run("func cancel", func(t *testing.T) {
		called := false

		timeout := deadline.New(time.Second)
		err := timeout.Run(func(cancel func()) error {
			called = true
			cancel()
			return errors.New("bad")
		})

		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})

	t.Run("func timeout", func(t *testing.T) {
		called := false

		timeout := deadline.New(100 * time.Millisecond)
		err := timeout.Run(func(cancel func()) error {
			called = true
			time.Sleep(time.Millisecond * 200)
			return errors.New("bad")
		})

		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, deadline.ErrTimeout(err); expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})
}

func Example() {
	timeout := deadline.New(10 * time.Millisecond)
	err := timeout.Run(func(cancel func()) error {
		return nil
	})

	switch {
	case err == nil:
		fmt.Println("success!")
	case deadline.ErrTimeout(err):
		fmt.Println("deadline timeout")
	default:
		fmt.Println("other error")
	}

	// Output:
	// success!
}
