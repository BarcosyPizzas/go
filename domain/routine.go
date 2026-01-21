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
	Exercise Exercise
	Sets     int
	Reps     int
}

func NewExerciseDetail(exercise Exercise, sets, reps int) ExerciseDetail {
	if sets == 0 {
		sets = 3
	}
	if reps == 0 {
		// low reps heavy weight yea buddyyyyy
		reps = 6
	}
	return ExerciseDetail{
		Exercise: exercise,
		Sets:     sets,
		Reps:     reps,
	}
}

// Need to think more about this, im troling, I hate cursor stop creating comments for me fakin shit.
func ExerciseDetails(exercises []Exercise) []ExerciseDetail {
	exerciseDetails := []ExerciseDetail{}
	for _, exercise := range exercises {
		exerciseDetails = append(exerciseDetails, NewExerciseDetail(exercise, 3, 6))
	}
	return exerciseDetails
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
