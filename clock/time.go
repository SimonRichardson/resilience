package clock

import (
	"sync"
	"time"
)

type timeClock struct {
	mutex sync.RWMutex
	stamp time.Time
}

// NewTimeClock creates a new time clock
func NewTimeClock() Clock {
	return &timeClock{
		mutex: sync.RWMutex{},
		stamp: time.Now(),
	}
}

// Time is used to return the current value of the clock
func (l *timeClock) Now() Time {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return timeTime(l.stamp.UnixNano())
}

// Increment is used to increment and return the value of the clock
func (l *timeClock) Increment() Time {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.stamp = l.stamp.Add(time.Nanosecond)
	return timeTime(l.stamp.UnixNano())
}

// Witness is called to update our local clock if necessary after
// witnessing a clock value received from another process
func (l *timeClock) Witness(v Time) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

}

type timeTime uint64

func (t timeTime) Value() uint64 {
	return uint64(t)
}
