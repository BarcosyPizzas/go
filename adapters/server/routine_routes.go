package server

import (
	"encoding/json"
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
