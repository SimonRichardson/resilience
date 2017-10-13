package timer

import (
	"sync"
	"time"
)

type timeTimer struct {
	mutex   sync.RWMutex
	current time.Time
	expiry  time.Duration
	timer   *time.Timer
}

func NewTimer(expiry time.Duration, fn func()) Timer {
	return &timeTimer{
		current: time.Now(),
		expiry:  expiry,
		timer:   time.AfterFunc(expiry, fn),
	}
}

func (t *timeTimer) Now() Time {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.current = time.Now()
	return timeTime(uint64(t.current.UnixNano()))
}

func (t *timeTimer) After() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	expiry := t.current.Add(t.expiry)
	return time.Now().After(expiry)
}

func (t *timeTimer) Reset() error {
	t.timer.Stop()
	t.timer.Reset(t.expiry)
	return nil
}

type timeTime uint64

func (t timeTime) Value() uint64 {
	return uint64(t)
}
