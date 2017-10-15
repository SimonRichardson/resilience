package gate_test

import (
	"fmt"
	"testing"

	"github.com/SimonRichardson/resilience/gate"
)

func Test(t *testing.T) {
	t.Parallel()

	t.Run("left func called", func(t *testing.T) {
		called := false

		branch := gate.New(
			func() error {
				called = true
				return nil
			},
			func() error {
				t.Fatal("failed if called")
				return nil
			},
		)
		err := branch.Run()

		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})

	t.Run("right func called", func(t *testing.T) {
		called := false

		branch := gate.New(
			func() error {
				t.Fatal("failed if called")
				return nil
			},
			func() error {
				called = true
				return nil
			},
		)
		branch.Switch()
		err := branch.Run()

		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})

	t.Run("calling switch multiple times func called", func(t *testing.T) {
		called := false

		branch := gate.New(
			func() error {
				t.Fatal("failed if called")
				return nil
			},
			func() error {
				called = true
				return nil
			},
		)
		branch.Switch()
		branch.Switch()
		branch.Switch()
		err := branch.Run()

		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})
}

func Example() {
	branch := gate.New(
		func() error {
			fmt.Println("left")
			return nil
		},
		func() error {
			fmt.Println("right")
			return nil
		},
	)

	err := branch.Run()

	switch {
	case err == nil:
		fmt.Println("success!")
	default:
		fmt.Println("other error")
	}

	branch.Switch()

	err = branch.Run()

	switch {
	case err == nil:
		fmt.Println("success!")
	default:
		fmt.Println("other error")
	}

	// Output:
	// left
	// success!
	// right
	// success!
}
