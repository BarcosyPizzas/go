package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	servermux := http.NewServeMux()
	servermux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	servermux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Home Page")
	})
	servermux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "About Page")
	})
	servermux.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		var humbe humbe
		if err := json.NewDecoder(r.Body).Decode(&humbe); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode body because you are a little piece of shit: %v", err), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Beautiful song I love it: %s", humbe.Title)
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: servermux,
	}
	log.Fatal(server.ListenAndServe())
}

// humbe is my struct to store body of song i like
type humbe struct {
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Artist   string `json:"artist"`
}

type database struct {
	db *sql.DB
}
