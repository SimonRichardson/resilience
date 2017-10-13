package deadline

import (
	"time"
)

type Deadline struct {
	timeout time.Duration
}

func New(timeout time.Duration) *Deadline {
	return &Deadline{
		timeout: timeout,
	}
}

func (d *Deadline) Run(fn func(func()) error) error {
	var (
		err = make(chan error)
		res = make(chan struct{})
	)

	go func() {
		err <- fn(func() {
			res <- struct{}{}
		})
	}()

	select {
	case <-res:
		return nil
	case err := <-err:
		return err
	case <-time.After(d.timeout):
		close(res)
		return errTimeout{}
	}
}

type errTimeout struct{}

func (e errTimeout) Error() string {
	return "deadline timed out"
}

// ErrTimeout checks if the error was because of deadline timedout.
func ErrTimeout(err error) bool {
	_, ok := err.(errTimeout)
	return ok
}
