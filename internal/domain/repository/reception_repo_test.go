package repository_test

import (
	"pvz/internal/domain/model"
	"pvz/internal/domain/repository"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateReception(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()
	repo := repository.NewReceptionRepo(db)
	pvzID := uuid.New()
	reception := &model.Reception{
		PVZID:  pvzID,
		Status: "in_progress",
	}
	reception.ID = uuid.New()
	reception.DateTime = time.Now()

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO receptions (id, date_time, pvz_id, status)
		VALUES ($1, $2, $3, $4)
	`)).
		WithArgs(reception.ID, reception.DateTime, reception.PVZID, reception.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.CreateReception(reception)
	assert.NoError(t, err)
}

func TestHasOpenReception(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()
	repo := repository.NewReceptionRepo(db)
	pvzID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*) FROM receptions 
		WHERE pvz_id = $1 AND status = 'in_progress'
	`)).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	hasOpen, err := repo.HasOpenReception(pvzID)
	assert.NoError(t, err)
	assert.True(t, hasOpen)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLastOpenReceptionByPVZ(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()
	repo := repository.NewReceptionRepo(db)
	pvzID := uuid.New()
	expectedReception := &model.Reception{
		ID:       uuid.New(),
		DateTime: time.Now(),
		PVZID:    pvzID,
		Status:   "in_progress",
	}
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, date_time, pvz_id, status 
		FROM receptions 
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1
	`)).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
			AddRow(expectedReception.ID, expectedReception.DateTime, expectedReception.PVZID, expectedReception.Status))
	reception, err := repo.FindLastOpenReceptionByPVZ(pvzID)
	require.NoError(t, err)
	require.NotNil(t, reception)
	assert.Equal(t, expectedReception.PVZID, reception.PVZID)
	assert.Equal(t, expectedReception.Status, reception.Status)
	assert.Equal(t, expectedReception.ID, reception.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLastOpenReceptionByPVZ_NotFound(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()
	repo := repository.NewReceptionRepo(db)

	pvzID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, date_time, pvz_id, status 
		FROM receptions 
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1
	`)).
		WithArgs(pvzID).
		WillReturnError(gorm.ErrRecordNotFound)

	reception, err := repo.FindLastOpenReceptionByPVZ(pvzID)
	assert.Error(t, err)
	assert.Nil(t, reception)
	assert.EqualError(t, err, "record not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseReception(t *testing.T) {
	db, mock, cleanup := SetupTestDB(t)
	defer cleanup()
	repo := repository.NewReceptionRepo(db)
	pvzID := uuid.New()
	receptionID := uuid.New()
	reception := &model.Reception{
		ID:     receptionID,
		PVZID:  pvzID,
		Status: "in_progress",
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE receptions SET status = 'close' WHERE id = $1`)).
		WithArgs(reception.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err := repo.CloseReception(reception)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
