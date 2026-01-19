package storage

import "gymlog/domain"

// Storage is the interface for the storage layer.
type Storage interface {
	Close() error
	Exercises() ([]domain.Exercise, error)
}
