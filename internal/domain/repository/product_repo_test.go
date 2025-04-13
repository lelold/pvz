package repository_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"pvz/internal/domain/model"
	"pvz/internal/domain/repository"
)

func TestCreateProduct(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewProductRepo(db)

	product := &model.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        "электроника",
		ReceptionID: uuid.New(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "products"`)).
		WithArgs(product.ID, product.DateTime, product.Type, product.ReceptionID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(product)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLastAddedProduct(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewProductRepo(db)

	receptionID := uuid.New()
	productID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE reception_id = $1 ORDER BY date_time DESC,"products"."id" LIMIT $2`)).
		WithArgs(receptionID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id"}).
			AddRow(productID, now, "электроника", receptionID))

	product, err := repo.GetLastAddedProduct(receptionID.String())
	assert.NoError(t, err)
	assert.Equal(t, "электроника", product.Type)
	assert.Equal(t, receptionID, product.ReceptionID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProductByID(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewProductRepo(db)

	productID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "products" WHERE id = $1`)).
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteProductByID(productID.String())
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
