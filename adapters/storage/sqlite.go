package storage

import (
	"database/sql"
	"gymlog/domain"
)

type sqliteStorage struct {
	db *sql.DB
}

// NewSqliteStorage creates a new SQLite storage.
func NewSqliteStorage(db *sql.DB) Storage {
	return &sqliteStorage{db: db}
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
