package application

import "gymlog/domain"

// RoutineRepository is the interface for the routine repository.
type RoutineRepository interface {
	Exercises() ([]domain.Exercise, error)
	SetRoutine(routine domain.Routine) error
}

type UserRepository interface {
	Users(username string) ([]domain.User, error)
	SaveUser(user domain.User) error
}
