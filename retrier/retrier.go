package retrier

import (
	"time"
)

type Retrier struct {
	backoff []time.Duration
}

func New(amount int, duration time.Duration) *Retrier {
	return &Retrier{
		backoff: exponential(amount, duration),
	}
}

// Run executes the given function
func (r *Retrier) Run(fn func() error) error {
	var retries int
	for {
		err := fn()
		if err == nil {
			return nil
		}

		if retries >= len(r.backoff) {
			return errRetry{err}
		}
		time.Sleep(r.backoff[retries])
		retries++
	}
}

func exponential(n int, d time.Duration) []time.Duration {
	res := make([]time.Duration, n)
	for i := 0; i < n; i++ {
		res[i] = d
		d *= 2
	}
	return res
}

type errRetry struct {
	err error
}

func (e errRetry) Error() string {
	return e.err.Error()
}

// ErrRetry checks if the error was because of too many retries.
func ErrRetry(err error) bool {
	_, ok := err.(errRetry)
	return ok
}
