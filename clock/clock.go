package clock

// Clock is a thread safe implementation of a lamport clock.
type Clock interface {
	Now() Time
	Increment() Time
	Witness(Time)
}

// Time is the value of a Clock.
type Time interface {
	Value() uint64
}
