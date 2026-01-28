package storage

import (
	"database/sql"
	"gymlog/domain"

	_ "github.com/mattn/go-sqlite3"
)

// sqliteStorage is the implementation of the Storage interface for SQLite.
type sqliteStorage struct {
	db *sql.DB
}

// NewSqliteStorage creates a new SQLite storage.
func NewSqliteStorage(dbPath string) (Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &sqliteStorage{db: db}, nil
}

func (s *sqliteStorage) Close() error {
	return s.db.Close()
}

// Exercises returns all the exercises from the database.
func (s *sqliteStorage) Exercises() ([]domain.Exercise, error) {
	rows, err := s.db.Query("SELECT id, name, target FROM exercises")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := []domain.Exercise{}
	for rows.Next() {
		var exercise domain.Exercise
		err = rows.Scan(&exercise.ID, &exercise.Name, &exercise.Target)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, nil
}

func (s *sqliteStorage) SaveRoutine(userID int, routine domain.Routine) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO routines (name, description, user_id) VALUES (?, ?, ?)", routine.Name, routine.Description, userID)
	if err != nil {
		return err
	}
	routineID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for i, exercise := range routine.Exercises {
		_, err = tx.Exec("INSERT INTO routine_exercises (routine_id, exercise_id, order_index, sets, reps) VALUES (?, ?, ?, ?, ?)", routineID, exercise.ID, i, exercise.Sets, exercise.Reps)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *sqliteStorage) Users(username string) ([]domain.User, error) {
	rows, err := s.db.Query("SELECT id, username, email, password_hash FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *sqliteStorage) SaveUser(username string, email string, passwordHash string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", username, email, passwordHash)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqliteStorage) SaveSession(userID int, sessionToken string, csrfToken string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("INSERT INTO sessions (user_id, session_token, csrf_token) VALUES (?, ?, ?)", userID, sessionToken, csrfToken)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqliteStorage) GetUserSession(userID int) (domain.UserSession, error) {
	row := s.db.QueryRow("SELECT session_token, csrf_token FROM sessions WHERE user_id = ?", userID)

	var domainUserSession domain.UserSession
	err := row.Scan(&domainUserSession.SessionToken, &domainUserSession.CSRFToken)
	if err != nil {
		return domain.UserSession{}, err
	}
	domainUserSession.UserID = userID
	return domainUserSession, nil
}

func (s *sqliteStorage) DeleteSession(userID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
