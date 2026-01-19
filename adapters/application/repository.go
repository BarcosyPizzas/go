package application

import "gymlog/domain"

type RoutineRepository interface {
	Exercises() ([]domain.Exercise, error)
}
