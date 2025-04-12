package repository_test

import (
	"pvz/internal/domain/model"
	"pvz/internal/domain/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}
	return gormDB, mock, func() {
		db.Close()
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	userID := uuid.New()
	user := &model.User{
		ID:    userID,
		Email: "test@example.com",
		Role:  "client",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(user.ID, user.Email, user.Password, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_Found(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewUserRepository(db)
	email := "test@example.com"
	userID := uuid.New().String()
	rows := sqlmock.NewRows([]string{"id", "email", "password", "role"}).
		AddRow(userID, email, "123", "employee")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)
	user, err := repo.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewUserRepository(db)
	email := "test@example.com"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnError(gorm.ErrRecordNotFound)
	user, err := repo.GetUserByEmail(email)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
