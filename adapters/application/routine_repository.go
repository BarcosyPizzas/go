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

func (r *GymRepository) SetRoutine(userID int, routine domain.Routine) error {
	if len(routine.Exercises) == 0 {
		return errors.New("routine must have at least one exercise")
	}
	return r.storage.SaveRoutine(userID, routine)
}

func (r *GymRepository) GetRoutines(userID int) ([]domain.Routine, error) {
	routines, err := r.storage.Routines(userID)
	if err != nil {
		return nil, err
	}
	return routines, nil
}

func (r *GymRepository) GetRoutine(routineID int) (domain.Routine, error) {
	routine, err := r.storage.Routine(routineID)
	if err != nil {
		return domain.Routine{}, err
	}
	return routine, nil
}
