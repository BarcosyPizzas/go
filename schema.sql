-- Tabla de usuarios
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de ejercicios (compartidos entre todos los usuarios)
CREATE TABLE exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    target VARCHAR(255) NOT NULL
);

-- Tabla de rutinas (ahora asociadas a usuarios)
CREATE TABLE routines (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Tabla de unión para la relación muchos-a-muchos entre rutinas y ejercicios
CREATE TABLE routine_exercises (
    routine_id INTEGER NOT NULL,
    exercise_id INTEGER NOT NULL,
    order_index INTEGER NOT NULL, -- Para mantener el orden de los ejercicios en la rutina
    sets INTEGER, -- Número de series para este ejercicio en esta rutina
    reps INTEGER, -- Número de repeticiones para este ejercicio en esta rutina
    PRIMARY KEY (routine_id, exercise_id),
    FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

-- Índices para mejorar el rendimiento de las consultas
CREATE INDEX idx_routines_user_id ON routines(user_id);
CREATE INDEX idx_routine_exercises_routine_id ON routine_exercises(routine_id);
CREATE INDEX idx_routine_exercises_exercise_id ON routine_exercises(exercise_id);
CREATE INDEX idx_routine_exercises_order ON routine_exercises(routine_id, order_index);