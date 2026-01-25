package server

import (
	"gymlog/adapters/application"
	"log"
	"net/http"
)

// server is the entry point for http requests for swift/nextjs frontend.
type gymlogServer struct {
	server            *http.Server
	routineRepository application.RoutineRepository
	userRepository    application.UserRepository
}

// NewServer is the constructor for the server.
func NewServer(routineRepository application.RoutineRepository, userRepository application.UserRepository) *gymlogServer {

	s := &gymlogServer{
		routineRepository: routineRepository,
		userRepository:    userRepository,
	}

	s.server = &http.Server{
		Addr:    ":6767",
		Handler: s.loadHandlers(),
	}
	return s
}

// loadHandlers loads all the handlers for the server.
func (s *gymlogServer) loadHandlers() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})
	handler.HandleFunc("/exercises", s.handleGetExercises)
	handler.HandleFunc("/register", s.handleRegister)
	return handler
}

// Start simply starts the server.
func (s *gymlogServer) Start() error {
	log.Println("Starting server on port 6767")
	return s.server.ListenAndServe()
}
