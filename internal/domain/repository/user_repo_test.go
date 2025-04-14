package repository_test

import (
	"database/sql"
	"pvz/internal/domain/model"
	"pvz/internal/domain/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}

func TestCreateUser(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	userID := uuid.New()
	user := &model.User{
		ID:    userID,
		Email: "test@example.com",
		Role:  "client",
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO users (id, email, password, role)
		VALUES ($1, $2, $3, $4)
	`)).
		WithArgs(user.ID, user.Email, user.Password, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_Found(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()
	repo := repository.NewUserRepository(db)
	email := "test@example.com"
	userID := uuid.New().String()
	rows := sqlmock.NewRows([]string{"id", "email", "password", "role"}).
		AddRow(userID, email, "123", "employee")
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, email, password, role
		FROM users
		WHERE email = $1
	`)).
		WithArgs(email).
		WillReturnRows(rows)
	user, err := repo.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()
	repo := repository.NewUserRepository(db)
	email := "test@example.com"
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, email, password, role
		FROM users
		WHERE email = $1
	`)).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)
	user, err := repo.GetUserByEmail(email)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
