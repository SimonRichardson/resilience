package clock

import (
	"sync/atomic"
	"time"
)

type timeClock struct {
	stamp uint64
}

// NewTimeClock creates a new time clock
func NewTimeClock() Clock {
	return &timeClock{
		stamp: uint64(time.Now().UnixNano()),
	}
}

// Time is used to return the current value of the clock
func (l *timeClock) Now() Time {
	return timeTime(atomic.LoadUint64(&l.stamp))
}

// Increment is used to increment and return the value of the clock
func (l *timeClock) Increment() Time {
	return timeTime(atomic.AddUint64(&l.stamp, 1))
}

// Witness is called to update our local clock if necessary after
// witnessing a clock value received from another process
func (l *timeClock) Witness(v Time) {
WITNESS:
	// if the other value is old, we do not need to do anything
	var (
		cur   = atomic.LoadUint64(&l.stamp)
		other = v.Value()
	)
	if other < cur {
		return
	}

	// Ensure that our local clock is at least one ahead.
	if !atomic.CompareAndSwapUint64(&l.stamp, cur, other+1) {
		// The CAS failed, so we just retry. Eventually our CAS should
		// succeed or a future witness will pass us by and our witness
		// will end.
		goto WITNESS
	}
}

func (l *timeClock) Clone() Clock {
	return &timeClock{stamp: atomic.LoadUint64(&l.stamp)}
}

func (l *timeClock) Reset() {
RESET:
	var (
		cur   = atomic.LoadUint64(&l.stamp)
		other = uint64(time.Now().UnixNano())
	)
	if other == cur {
		return
	}

	if !atomic.CompareAndSwapUint64(&l.stamp, cur, other) {
		goto RESET
	}
}

type timeTime uint64

func (t timeTime) Value() uint64 {
	return uint64(t)
}

func (t timeTime) Before(other Time) bool {
	return uint64(t) < uint64(other.Value())
}

func (t timeTime) After(other Time) bool {
	return uint64(t) > uint64(other.Value())
}
