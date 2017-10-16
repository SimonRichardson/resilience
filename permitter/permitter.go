package permitter

// Permitter defines if a given action is permitted or not
type Permitter interface {

	// Run checks if a permitt is valid or not then runs the function
	Run(func() error) error
}
