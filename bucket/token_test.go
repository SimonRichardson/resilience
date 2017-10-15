package bucket_test

import (
	"fmt"
	"testing"

	"github.com/SimonRichardson/resilience/bucket"
)

func TestTokenBucket(t *testing.T) {
	t.Parallel()

	t.Run("take", func(t *testing.T) {
		tokens := bucket.NewTokenBucket(2)

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("multiple takes", func(t *testing.T) {
		tokens := bucket.NewTokenBucket(2)

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := int64(1), tokens.Take(2); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("put", func(t *testing.T) {
		tokens := bucket.NewTokenBucket(2)

		if expected, actual := int64(0), tokens.Put(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("take then put", func(t *testing.T) {
		tokens := bucket.NewTokenBucket(2)

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}

		if expected, actual := int64(1), tokens.Put(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})
}

func Example() {
	tokens := bucket.NewTokenBucket(2)
	fmt.Println(tokens.Take(1))
	fmt.Println(tokens.Take(1))
	fmt.Println(tokens.Take(0))
	fmt.Println(tokens.Put(2))
	fmt.Println(tokens.Put(1))

	// Output:
	// 1
	// 1
	// 0
	// 2
	// 0
}
