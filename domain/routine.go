package domain

// routine defines a list of exercises that compose a workout, for example push day.
type Routine struct {
	ID          int
	Name        string
	Description string
	Exercises   []Exercise
}
