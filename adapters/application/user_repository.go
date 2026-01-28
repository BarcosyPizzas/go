package application

import (
	"errors"
	"gymlog/adapters/storage"
	"gymlog/domain"
)

type UserRepo struct {
	storage storage.Storage
}

func NewUserRepo(storage storage.Storage) UserRepository {
	return &UserRepo{storage: storage}
}

func (r *UserRepo) Users(username string) ([]domain.User, error) {
	users, err := r.storage.Users(username)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) UserSession(username string) (domain.UserSession, error) {
	user, err := r.Users(username)
	if err != nil {
		return domain.UserSession{}, err
	}
	if len(user) == 0 {
		return domain.UserSession{}, errors.New("user not found")
	}
	return r.storage.GetUserSession(user[0].ID)
}

func (r *UserRepo) SaveUser(user domain.User) error {
	return r.storage.SaveUser(user.Username, user.Email, user.PasswordHash)
}

func (r *UserRepo) SaveSession(userID int, sessionToken string, csrfToken string) error {
	return r.storage.SaveSession(userID, sessionToken, csrfToken)
}

func (r *UserRepo) DeleteSession(userID int) error {
	return r.storage.DeleteSession(userID)
}
