package domain

// user defines a user of the application.
type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
}

type UserSession struct {
	UserID       int
	SessionToken string
	CSRFToken    string
}
