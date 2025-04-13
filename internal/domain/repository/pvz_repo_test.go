package repository_test

import (
	"regexp"
	"testing"
	"time"

	"pvz/internal/domain/model"
	"pvz/internal/domain/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatePVZ(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPVZRepo(db)

	pvz := &model.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Москва",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "pvzs"`)).
		WithArgs(pvz.ID, pvz.RegistrationDate, pvz.City).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(pvz)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFilteredPVZs(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPVZRepo(db)

	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()
	page := 1
	limit := 2

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "pvzs" WHERE registration_date >= $1 AND registration_date <= $2 LIMIT $3`)).
		WithArgs(start, end, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow(uuid.New(), time.Now(), "Москва").
			AddRow(uuid.New(), time.Now(), "Санкт-Петербург"))

	pvzs, err := repo.GetFilteredPVZs(&start, &end, page, limit)
	assert.NoError(t, err)
	assert.Len(t, pvzs, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetReceptionsWithProducts(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewPVZRepo(db)

	pvzID := uuid.New()
	receptionID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "receptions" WHERE pvz_id = $1`)).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "status", "date_time"}).
			AddRow(receptionID, pvzID, "in_progress", time.Now()))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE reception_id = $1`)).
		WithArgs(receptionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "reception_id", "name", "quantity"}).
			AddRow(uuid.New(), receptionID, "Товар 1", 3).
			AddRow(uuid.New(), receptionID, "Товар 2", 1))

	result, err := repo.GetReceptionsWithProducts(pvzID.String())
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Len(t, result[0].Products, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
