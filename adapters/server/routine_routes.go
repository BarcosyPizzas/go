package server

import (
	"encoding/json"
	"gymlog/domain"
	"net/http"
)

// handleGetExercises handles the GET request for the exercises.
func (s *gymlogServer) handleGetExercises(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Must be a GET request", http.StatusMethodNotAllowed)
		return
	}

	exercises, err := s.routineRepository.Exercises()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(exercises)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleSetRoutine sets a routine for a user.
func (s *gymlogServer) handleSetRoutine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Must be a POST request", http.StatusMethodNotAllowed)
		return
	}

	if err := s.Authorize(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	user, err := s.userRepository.Users(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(user) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var routineRequest postRoutineRequest
	err = json.NewDecoder(r.Body).Decode(&routineRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exerciseDetails := routineRequestToExerciseDetails(routineRequest)

	routine, err := domain.CreateRoutine(routineRequest.Name, routineRequest.Description, exerciseDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.routineRepository.SetRoutine(user[0].ID, routine)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

type postRoutineRequest struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Exercises   []postRoutineExercise `json:"exercises"`
}

type postRoutineExercise struct {
	ID   int `json:"id"`
	Sets int `json:"sets"`
	Reps int `json:"reps"`
}

func routineRequestToExerciseDetails(request postRoutineRequest) []domain.ExerciseDetail {
	exerciseDetails := []domain.ExerciseDetail{}
	for _, exercise := range request.Exercises {
		exerciseDetails = append(exerciseDetails, domain.NewExerciseDetail(exercise.ID, exercise.Sets, exercise.Reps))
	}
	return exerciseDetails
}
