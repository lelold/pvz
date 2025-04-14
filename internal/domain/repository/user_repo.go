package repository

import (
	"database/sql"
	"errors"
	"pvz/internal/domain/model"

	"github.com/google/uuid"
)

type UserRepo interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *model.User) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	query := `
		INSERT INTO users (id, email, password, role)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(query, user.ID, user.Email, user.Password, user.Role)
	return err
}

func (r *userRepo) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, email, password, role
		FROM users
		WHERE email = $1
	`
	var user model.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
