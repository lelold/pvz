package repository

import (
	"errors"
	"pvz/internal/domain/model"
	"sync"
)

type UserRepo struct {
	mu    sync.RWMutex
	users map[string]model.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{users: make(map[string]model.User)}
}

func (r *UserRepo) Save(user model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Email] = user
	return nil
}

func (r *UserRepo) GetByEmail(email string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[email]
	if !exists {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}
