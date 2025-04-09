package service

import (
	"errors"
	"pvz/internal/domain/model"
	repository "pvz/internal/repository/pg"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, password, role string) (model.User, error) {
	if role != "employee" && role != "moderator" {
		return model.User{}, errors.New("invalid role")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hashed),
		Role:     role,
	}
	err = s.repo.Save(user)
	return user, err
}

func (s *UserService) Login(email, password string) (model.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, errors.New("invalid credentials")
	}
	return user, nil
}
