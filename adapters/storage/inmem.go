package storage

import "gymlog/domain"

type inmemStorage struct {
	exercises map[int]domain.Exercise
	routines  map[int]domain.Routine
}

func NewInmemStorage() Storage {
	return &inmemStorage{
		exercises: make(map[int]domain.Exercise),
		routines:  make(map[int]domain.Routine),
	}
}

func (s *inmemStorage) SaveRoutine(routine domain.Routine) error {
	s.routines[routine.ID] = routine
	return nil
}

// Close closes the storage.
func (s *inmemStorage) Close() error {
	return nil
}

// Exercises returns all the exercises from the storage.
func (s *inmemStorage) Exercises() ([]domain.Exercise, error) {
	return []domain.Exercise{
		{ID: 1, Name: "Push-up", Target: "Chest"},
		{ID: 2, Name: "Pull-up", Target: "Back"},
		{ID: 3, Name: "Squat", Target: "Legs"},
	}, nil
}
