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
	db, mock, cleanup := setupTestDB(t)
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
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "receptions"`)).
		WithArgs(reception.ID, reception.DateTime, reception.PVZID, reception.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := repo.CreateReception(reception)
	assert.NoError(t, err)
}

func TestHasOpenReception(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewReceptionRepo(db)
	pvzID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "receptions"`)).
		WithArgs(pvzID, "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	hasOpen, err := repo.HasOpenReception(pvzID.String())
	assert.NoError(t, err)
	assert.True(t, hasOpen)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLastOpenReceptionByPVZ(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewReceptionRepo(db)
	pvzID := uuid.New()
	expectedReception := &model.Reception{
		ID:       uuid.New(),
		DateTime: time.Now(),
		PVZID:    pvzID,
		Status:   "in_progress",
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "receptions" WHERE pvz_id = $1 AND status = $2 ORDER BY date_time DESC,"receptions"."id" LIMIT $3`)).
		WithArgs(pvzID, "in_progress", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
			AddRow(expectedReception.ID, expectedReception.DateTime, expectedReception.PVZID, expectedReception.Status))
	reception, err := repo.FindLastOpenReceptionByPVZ(pvzID.String())
	require.NoError(t, err)
	require.NotNil(t, reception)
	assert.Equal(t, expectedReception.PVZID, reception.PVZID)
	assert.Equal(t, expectedReception.Status, reception.Status)
	assert.Equal(t, expectedReception.ID, reception.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLastOpenReceptionByPVZ_NotFound(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repository.NewReceptionRepo(db)

	pvzID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "receptions" WHERE pvz_id = $1 AND status = $2 ORDER BY date_time DESC,"receptions"."id" LIMIT $3`)).
		WithArgs(pvzID, "in_progress", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	reception, err := repo.FindLastOpenReceptionByPVZ(pvzID.String())
	assert.Error(t, err)
	assert.Nil(t, reception)
	assert.EqualError(t, err, "открытая приемка не найдена")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseReception(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
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
		`UPDATE "receptions" SET "date_time"=$1,"pvz_id"=$2,"status"=$3 WHERE "id" = $4`)).
		WithArgs(reception.DateTime, reception.PVZID, "close", reception.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err := repo.CloseReception(reception)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
