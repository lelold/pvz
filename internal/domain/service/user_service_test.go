package service_test

import (
	"errors"
	"pvz/internal/domain/mocks"
	"pvz/internal/domain/model"
	"pvz/internal/domain/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	mockRepo := new(mocks.UserRepo)
	svc := service.NewUserService(mockRepo)

	t.Run("успешная регистрация", func(t *testing.T) {
		email := "test@example.com"
		password := "123"
		role := "employee"

		mockRepo.On("CreateUser", mock.MatchedBy(func(u *model.User) bool {
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
			return u.Email == email && u.Role == role && err == nil
		})).Return(nil).Once()

		user, err := svc.Register(email, password, role)

		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, role, user.Role)
		assert.NotEmpty(t, user.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("неверная роль", func(t *testing.T) {
		_, err := svc.Register("test@example.com", "pass", "client")
		assert.Error(t, err)
		assert.Equal(t, "неверная роль", err.Error())
	})

	t.Run("ошибка хеширования", func(t *testing.T) {
		longPassword := make([]byte, 1000000)
		_, err := svc.Register("test@example.com", string(longPassword), "employee")
		assert.Error(t, err)
	})
}

func TestLogin(t *testing.T) {
	mockRepo := new(mocks.UserRepo)
	svc := service.NewUserService(mockRepo)

	email := "test@example.com"
	password := "123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashed),
		Role:     "employee",
	}

	t.Run("успешный вход", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", email).Return(user, nil).Once()

		result, err := svc.Login(email, password)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Role, result.Role)

		mockRepo.AssertExpectations(t)
	})

	t.Run("неверный пароль", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", email).Return(user, nil).Once()

		_, err := svc.Login(email, "wrongpass")
		assert.Error(t, err)
		assert.Equal(t, "неверные данные", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "nf@test.com").Return(nil, errors.New("not found")).Once()

		_, err := svc.Login("nf@test.com", "any")
		assert.Error(t, err)
	})
}
