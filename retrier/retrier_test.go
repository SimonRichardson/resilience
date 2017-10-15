package retrier_test

import (
	"fmt"
	"time"

	"github.com/SimonRichardson/resilience/retrier"
)

func Example() {
	retry := retrier.New(10, time.Second)
	err := retry.Run(func() error {
		return nil
	})

	switch {
	case err == nil:
		fmt.Println("success!")
	case retrier.ErrRetry(err):
		fmt.Println("deadline timeout")
	default:
		fmt.Println("other error")
	}

	// Output:
	// success!
}
