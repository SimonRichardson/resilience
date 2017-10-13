package breaker

import (
	"sync/atomic"
	"time"

	"github.com/SimonRichardson/resilience/clock"
	"github.com/SimonRichardson/resilience/timer"
)

const (
	openState uint32 = iota
	halfOpenState
	closedState
)

type CircuitBreaker struct {
	state      uint32
	metrics    *metrics
	thresholds *thresholds
	timer      timer.Timer
}

func New(success, failure uint64, expiry time.Duration) *CircuitBreaker {
	breaker := &CircuitBreaker{
		state:      closedState,
		metrics:    newMetrics(),
		thresholds: newThresholds(success, failure),
	}
	breaker.timer = timer.NewTimer(expiry, breaker.tick)
	return breaker
}

func (c *CircuitBreaker) Run(fn func() error) error {
	state := atomic.LoadUint32(&c.state)
	if state == openState {
		return errBreakerOpen{}
	}

	err := fn()

	if c.metrics.Failed.Now().Value() > 0 {
		if c.timer.After() {
			c.metrics.Reset()
		}
	}

	switch state {
	case closedState:
		if err == nil {
			break
		}

		val := c.metrics.Failed.Increment()
		if c.thresholds.AfterFailed(val.Value()) {
			c.open()
		} else {
			c.timer.Now()
		}

	case halfOpenState:
		if err != nil {
			c.open()
			break
		}

		val := c.metrics.Success.Increment()
		if c.thresholds.AfterSuccess(val.Value()) {
			c.close()
		}
	}

	return err
}

func (c *CircuitBreaker) open() {
	c.metrics.Reset()
	atomic.StoreUint32(&c.state, openState)

	c.timer.Reset()
}

func (c *CircuitBreaker) close() {
	c.metrics.Reset()
	atomic.StoreUint32(&c.state, closedState)
}

func (c *CircuitBreaker) tick() {
	if state := atomic.LoadUint32(&c.state); state == closedState {
		return
	}

	c.metrics.Reset()
	atomic.StoreUint32(&c.state, halfOpenState)
}

type metrics struct {
	Success, Failed clock.Clock
}

func newMetrics() *metrics {
	return &metrics{
		Success: clock.NewLamportClock(),
		Failed:  clock.NewLamportClock(),
	}
}

func (m metrics) Reset() {
	m.Success.Reset()
	m.Failed.Reset()
}

type thresholds struct {
	Success, Failed uint64
}

func newThresholds(success, failure uint64) *thresholds {
	return &thresholds{
		Success: success,
		Failed:  failure,
	}
}

func (t *thresholds) AfterSuccess(x uint64) bool {
	return x >= t.Success
}

func (t *thresholds) AfterFailed(x uint64) bool {
	return x >= t.Failed
}

type errBreakerOpen struct{}

func (e errBreakerOpen) Error() string {
	return "circuit breaker is open"
}

// ErrBreakerOpen checks if the error was because of the breaker being in the
// open state.
func ErrBreakerOpen(err error) bool {
	_, ok := err.(errBreakerOpen)
	return ok
}
