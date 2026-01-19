-- Tabla de ejercicios
CREATE TABLE exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Tabla de rutinas
CREATE TABLE routines (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT
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
CREATE INDEX idx_routine_exercises_routine_id ON routine_exercises(routine_id);
CREATE INDEX idx_routine_exercises_exercise_id ON routine_exercises(exercise_id);
CREATE INDEX idx_routine_exercises_order ON routine_exercises(routine_id, order_index);

