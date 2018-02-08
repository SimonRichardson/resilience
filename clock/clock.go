package clock

// Clock is a thread safe implementation of a clock.
type Clock interface {

	// Time is used to return the current value of the clock
	Now() Time

	// Increment is used to increment and return the value of the clock
	Increment() Time

	// Witness is called to update our local clock if necessary after
	// witnessing a clock value received from another process
	Witness(Time)

	// Clone a clock with the same local time underneath
	Clone() Clock

	// Reset the clock
	Reset()
}

// Time is the value of a Clock.
type Time interface {
	Value() uint64
}
