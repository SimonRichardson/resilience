package timer

// Timer is reasonable implementation of a timer.
type Timer interface {

	// Time is used to return the current value of the timer
	Now() Time

	// After is used to know if the time now is after timer has ticked
	After() bool

	// Reset the timer
	Reset() error
}

// Time is the value of a Timer.
type Time interface {
	Value() uint64
}
