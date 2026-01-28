package application

import "gymlog/domain"

// RoutineRepository is the interface for the routine repository.
type RoutineRepository interface {
	Exercises() ([]domain.Exercise, error)
	SetRoutine(userID int, routine domain.Routine) error
}

type UserRepository interface {
	Users(username string) ([]domain.User, error)
	SaveUser(user domain.User) error
	SaveSession(userID int, sessionToken string, csrfToken string) error
	DeleteSession(userID int) error
	UserSession(username string) (domain.UserSession, error)
}
