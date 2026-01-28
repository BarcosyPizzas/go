package storage

import "gymlog/domain"

// Storage is the interface for the storage layer.
type Storage interface {
	Close() error
	Exercises() ([]domain.Exercise, error)
	SaveRoutine(userID int, routine domain.Routine) error
	Users(username string) ([]domain.User, error)
	SaveUser(username string, email string, passwordHash string) error
	SaveSession(userID int, sessionToken string, csrfToken string) error
	GetUserSession(userID int) (domain.UserSession, error)
	DeleteSession(userID int) error
}
