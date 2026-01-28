package server

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"gymlog/domain"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *gymlogServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Must be a POST request", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	users, err := s.userRepository.Users(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(users) > 0 {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepository.SaveUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}

func (s *gymlogServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Must be a POST request", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := s.userRepository.Users(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(user) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !checkPasswordHash(password, user[0].PasswordHash) {
		http.Error(w, "Incorrect password, never try again", http.StatusUnauthorized)
		return
	}

	sessionToken, err := generateToken(32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	csrfToken, err := generateToken(32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	// Store token in database
	err = s.userRepository.SaveSession(user[0].ID, sessionToken, csrfToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// handleLogout logs out a user by deleting the session from the database.
func (s *gymlogServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Must be a POST request", http.StatusMethodNotAllowed)
		return
	}

	if err := s.Authorize(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	user, err := s.userRepository.UserSession(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.SessionToken == "" || user.CSRFToken == "" {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: false,
	})

	err = s.userRepository.DeleteSession(user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

// handleGetSession returns the session token and CSRF token for a user.
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// checkPasswordHash checks if a password matches a hash.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateToken generates a random token of a given length.
// Usa RawURLEncoding (sin padding '=') para que los tokens funcionen bien
// en cookies, headers y al copiar/pegar en Postman u otros clientes.
func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (s *gymlogServer) Authorize(r *http.Request) error {
	username := r.FormValue("username")
	user, err := s.userRepository.UserSession(username)
	if err != nil {
		return err
	}
	if user.SessionToken == "" || user.CSRFToken == "" {
		return errors.New("Unauthorized")
	}

	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		return errors.New("Unauthorized")
	}

	csrf := r.Header.Get("X-CSRF-Token")
	if csrf == "" || csrf != user.CSRFToken {
		return errors.New("Unauthorized")
	}

	return nil
}
