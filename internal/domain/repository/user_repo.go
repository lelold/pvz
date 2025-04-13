package repository

import (
	"pvz/internal/domain/model"

	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
