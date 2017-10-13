package clock

import "sync/atomic"

type lamportClock struct {
	counter uint64
}

// NewLamportClock creates a new lamport clock
func NewLamportClock() Clock {
	return &lamportClock{counter: 0}
}

// Time is used to return the current value of the clock
func (l *lamportClock) Now() Time {
	return lamportTime(atomic.LoadUint64(&l.counter))
}

// Increment is used to increment and return the value of the clock
func (l *lamportClock) Increment() Time {
	return lamportTime(atomic.AddUint64(&l.counter, 1))
}

// Witness is called to update our local clock if necessary after
// witnessing a clock value received from another process
func (l *lamportClock) Witness(v Time) {
WITNESS:
	// If the other value is old, we do not need to do anything
	var (
		cur   = atomic.LoadUint64(&l.counter)
		other = v.Value()
	)
	if other < cur {
		return
	}

	// Ensure that our local clock is at least one ahead.
	if !atomic.CompareAndSwapUint64(&l.counter, cur, other+1) {
		// The CAS failed, so we just retry. Eventually our CAS should
		// succeed or a future witness will pass us by and our witness
		// will end.
		goto WITNESS
	}
}

func (l *lamportClock) Reset() {
RESET:
	cur := atomic.LoadUint64(&l.counter)
	if cur == 0 {
		return
	}

	if !atomic.CompareAndSwapUint64(&l.counter, cur, 0) {
		goto RESET
	}
}

type lamportTime uint64

func (t lamportTime) Value() uint64 {
	return uint64(t)
}
