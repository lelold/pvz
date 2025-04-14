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
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := repository.NewPVZRepo(db)

	pvz := &model.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Москва",
	}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO pvzs (id, registration_date, city)
		VALUES ($1, $2, $3)
	`)).
		WithArgs(pvz.ID, pvz.RegistrationDate, pvz.City).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := repo.Create(pvz)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFilteredPVZs(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := repository.NewPVZRepo(db)

	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()
	page := 1
	limit := 2

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, registration_date, city FROM pvzs WHERE 1=1 AND registration_date >= $1 AND registration_date <= $2 ORDER BY registration_date DESC LIMIT $3 OFFSET $4`)).
		WithArgs(start, end, 2, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow(uuid.New(), time.Now(), "Москва").
			AddRow(uuid.New(), time.Now(), "Санкт-Петербург"))

	pvzs, err := repo.GetFilteredPVZs(&start, &end, page, limit)
	assert.NoError(t, err)
	assert.Len(t, pvzs, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetReceptionsWithProducts(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := repository.NewPVZRepo(db)

	pvzID := uuid.New()
	receptionID := uuid.New()
	now := time.Now()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1")).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
			AddRow(receptionID, now, pvzID, "in_progress"))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, date_time, type, reception_id FROM products WHERE reception_id = $1")).
		WithArgs(receptionID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"}).
			AddRow(uuid.New(), time.Now(), "Товар 1", receptionID).
			AddRow(uuid.New(), time.Now(), "Товар 2", receptionID))

	mock.ExpectCommit()

	result, err := repo.GetReceptionsWithProducts(pvzID)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Len(t, result[0].Products, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
