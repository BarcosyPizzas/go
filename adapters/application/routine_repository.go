package application

import (
	"errors"
	"gymlog/adapters/storage"
	"gymlog/domain"
)

// GymRepository is in charge of application business logic.
type GymRepository struct {
	storage storage.Storage
}

// NewGymRepository is the constructor for the GymRepository.
func NewGymRepository(storage storage.Storage) RoutineRepository {
	return &GymRepository{storage: storage}
}

// Exercises returns all the exercises from the storage.
func (r *GymRepository) Exercises() ([]domain.Exercise, error) {
	exercises, err := r.storage.Exercises()
	if err != nil {
		return nil, err
	}
	return exercises, nil
}

func (r *GymRepository) SetRoutine(routine domain.Routine) error {
	if len(routine.Exercises) == 0 {
		return errors.New("routine must have at least one exercise")
	}
	return r.storage.SaveRoutine(routine)
}
