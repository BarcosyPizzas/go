package domain

import "errors"

// routine defines a list of exercises that compose a workout, for example push day.
type Routine struct {
	ID          int
	Name        string
	Description string
	Exercises   []ExerciseDetail
}

// ExerciseDetail defines a single exercise with its sets and reps.
type ExerciseDetail struct {
	ID   int
	Sets int
	Reps int
}

func NewExerciseDetail(id int, sets, reps int) ExerciseDetail {
	if sets == 0 {
		sets = 3
	}
	if reps == 0 {
		// low reps heavy weight yea buddyyyyy
		reps = 6
	}
	return ExerciseDetail{
		ID:   id,
		Sets: sets,
		Reps: reps,
	}
}

func CreateRoutine(name, description string, exercises []ExerciseDetail) (Routine, error) {
	if name == "" {
		return Routine{}, errors.New("name is required")
	}
	if len(exercises) == 0 {
		return Routine{}, errors.New("at least one exercise is required")
	}
	return Routine{
		Name:        name,
		Description: description,
		Exercises:   exercises,
	}, nil
}
