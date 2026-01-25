package application

import (
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

func (r *UserRepo) SaveUser(user domain.User) error {
	return r.storage.SaveUser(user.Username, user.Email, user.PasswordHash)
}
