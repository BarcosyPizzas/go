package domain

// routine defines a list of exercises that compose a workout, for example push day.
type routine struct {
	ID          int
	Name        string
	Description string
	Exercises   []exercise
}

func (r *routine) addExercises(exercises []exercise) {
	r.Exercises = append(r.Exercises, exercises...)
}

func (r *routine) removeExercise(exercise exercise) {
	for i, e := range r.Exercises {
		if e.ID == exercise.ID {
			r.Exercises = append(r.Exercises[:i], r.Exercises[i+1:]...)
			break
		}
	}
}
