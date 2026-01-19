package server

import (
	"log"
	"net/http"
)

// server is the entry point for http requests for swift/nextjs frontend.
type gymlogServer struct {
	server *http.Server
}

// NewServer is the constructor for the server.
func NewServer() *gymlogServer {
	return &gymlogServer{
		server: &http.Server{
			Addr:    ":6767",
			Handler: loadHandlers(),
		},
	}
}

// loadHandlers loads all the handlers for the server.
func loadHandlers() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})
	return handler
}

// Start simply starts the server.
func (s *gymlogServer) Start() error {
	log.Println("Starting server on port 6767")
	return s.server.ListenAndServe()
}
