package permitter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SimonRichardson/resilience/permitter"
)

func TestPermitter(t *testing.T) {
	t.Parallel()

	t.Run("func called", func(t *testing.T) {
		called := false

		permit := permitter.New(1, time.Second)
		err := permit.Run(func() error {
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
	permit := permitter.New(1, 10*time.Millisecond)
	err := permit.Run(func() error {
		return nil
	})

	switch {
	case err == nil:
		fmt.Println("success!")
	default:
		fmt.Println("other error")
	}

	// Output:
	// success!
}
