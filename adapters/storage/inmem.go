package storage

import "gymlog/domain"

type inmemStorage struct {
	exercises map[int]domain.Exercise
	routines  map[int]domain.Routine
	users     map[string]domain.User
	sessions  map[int]domain.UserSession
}

func NewInmemStorage() Storage {
	return &inmemStorage{
		exercises: make(map[int]domain.Exercise),
		routines:  make(map[int]domain.Routine),
		users:     make(map[string]domain.User),
		sessions:  make(map[int]domain.UserSession),
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

func (s *inmemStorage) Users(username string) ([]domain.User, error) {
	return []domain.User{s.users[username]}, nil
}

func (s *inmemStorage) SaveUser(username string, email string, passwordHash string) error {
	s.users[username] = domain.User{ID: 1, Username: username, Email: email, PasswordHash: passwordHash}
	return nil
}

func (s *inmemStorage) SaveSession(userID int, sessionToken string, csrfToken string) error {
	s.sessions[userID] = domain.UserSession{UserID: userID, SessionToken: sessionToken, CSRFToken: csrfToken}
	return nil
}
