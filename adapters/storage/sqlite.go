package storage

import (
	"database/sql"
	_ "embed"
	"gymlog/domain"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed
var exercisesSeedSQL string

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
	storage := &sqliteStorage{db: db}
	if err := storage.seedExercises(); err != nil {
		return nil, err
	}
	return storage, nil
}

// seedExercises inserts default exercises if the table is empty.
func (s *sqliteStorage) seedExercises() error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM exercises").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Already seeded
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute each INSERT statement from the embedded SQL file
	statements := strings.Split(exercisesSeedSQL, "\n")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err = tx.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
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

	// Delete any existing sessions for this user to maintain only one session
	_, err = tx.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

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

func (s *sqliteStorage) Routines(userID int) ([]domain.Routine, error) {
	rows, err := s.db.Query(`
		SELECT r.id, r.name, r.description, re.exercise_id, re.sets, re.reps 
		FROM routines r 
		LEFT JOIN routine_exercises re ON r.id = re.routine_id 
		WHERE r.user_id = ? 
		ORDER BY r.id, re.order_index`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routineMap := make(map[int]*domain.Routine)
	var routineOrder []int

	for rows.Next() {
		var routineID int
		var name, description string
		var exerciseID, sets, reps sql.NullInt64

		err = rows.Scan(&routineID, &name, &description, &exerciseID, &sets, &reps)
		if err != nil {
			return nil, err
		}

		// Check if we've seen this routine before
		if _, exists := routineMap[routineID]; !exists {
			routineMap[routineID] = &domain.Routine{
				ID:          routineID,
				Name:        name,
				Description: description,
				Exercises:   []domain.ExerciseDetail{},
			}
			routineOrder = append(routineOrder, routineID)
		}

		// Add exercise if it exists (LEFT JOIN may return NULLs)
		if exerciseID.Valid {
			routineMap[routineID].Exercises = append(routineMap[routineID].Exercises, domain.ExerciseDetail{
				ID:   int(exerciseID.Int64),
				Sets: int(sets.Int64),
				Reps: int(reps.Int64),
			})
		}
	}

	// Build result slice maintaining order
	routines := make([]domain.Routine, 0, len(routineOrder))
	for _, id := range routineOrder {
		routines = append(routines, *routineMap[id])
	}
	return routines, nil
}

func (s *sqliteStorage) Routine(routineID int) (domain.Routine, error) {
	rows, err := s.db.Query(`
		SELECT r.id, r.name, r.description, re.exercise_id, re.sets, re.reps 
		FROM routines r 
		LEFT JOIN routine_exercises re ON r.id = re.routine_id 
		WHERE r.id = ? 
		ORDER BY re.order_index`, routineID)
	if err != nil {
		return domain.Routine{}, err
	}
	defer rows.Close()

	var routine domain.Routine
	found := false

	for rows.Next() {
		var id int
		var name, description string
		var exerciseID, sets, reps sql.NullInt64

		err = rows.Scan(&id, &name, &description, &exerciseID, &sets, &reps)
		if err != nil {
			return domain.Routine{}, err
		}

		if !found {
			routine = domain.Routine{
				ID:          id,
				Name:        name,
				Description: description,
				Exercises:   []domain.ExerciseDetail{},
			}
			found = true
		}

		if exerciseID.Valid {
			routine.Exercises = append(routine.Exercises, domain.ExerciseDetail{
				ID:   int(exerciseID.Int64),
				Sets: int(sets.Int64),
				Reps: int(reps.Int64),
			})
		}
	}

	if !found {
		return domain.Routine{}, sql.ErrNoRows
	}

	return routine, nil
}
