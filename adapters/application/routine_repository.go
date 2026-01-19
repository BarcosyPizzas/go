package application

import (
	"gymlog/adapters/storage"
	"gymlog/domain"
)

type GymRepository struct {
	storage storage.Storage
}

func NewGymRepository(storage storage.Storage) RoutineRepository {
	return &GymRepository{storage: storage}
}

func (r *GymRepository) Exercises() ([]domain.Exercise, error) {
	exercises, err := r.storage.Exercises()
	if err != nil {
		return nil, err
	}
	return exercises, nil
}
