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

func (s *sqliteStorage) SaveRoutine(routine domain.Routine) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO routines (name, description) VALUES (?, ?)", routine.Name, routine.Description)
	if err != nil {
		return err
	}

	for i, exercise := range routine.Exercises {
		_, err = tx.Exec("INSERT INTO routine_exercises (routine_id, exercise_id, order_index, sets, reps) VALUES (?, ?, ?, ?, ?)", routine.ID, exercise.ID, i, exercise.Sets, exercise.Reps)
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
