package application

import "gymlog/domain"

// RoutineRepository is the interface for the routine repository.
type RoutineRepository interface {
	Exercises() ([]domain.Exercise, error)
}
