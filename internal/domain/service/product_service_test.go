package service_test

import (
	"errors"
	"testing"

	"pvz/internal/domain/mocks"
	"pvz/internal/domain/model"
	"pvz/internal/domain/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct_Success(t *testing.T) {
	pRepo := new(mocks.ProductRepo)
	rRepo := new(mocks.ReceptionRepo)
	svc := service.NewProductService(pRepo, rRepo)

	pvzID := uuid.New()
	reception := &model.Reception{ID: uuid.New(), Status: "in_progress"}
	rRepo.On("FindLastOpenReceptionByPVZ", pvzID).Return(reception, nil)

	pRepo.On("Create", mock.AnythingOfType("*model.Product")).Return(nil)

	product, err := svc.CreateProduct("обувь", pvzID)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "обувь", product.Type)
	pRepo.AssertExpectations(t)
	rRepo.AssertExpectations(t)
}

func TestCreateProduct_InvalidType(t *testing.T) {
	svc := service.NewProductService(nil, nil)
	product, err := svc.CreateProduct("влорыар", uuid.New())
	assert.Nil(t, product)
	assert.EqualError(t, err, "неверный товар")
}

func TestCreateProduct_NoReception(t *testing.T) {
	rRepo := new(mocks.ReceptionRepo)
	svc := service.NewProductService(nil, rRepo)

	pvzID := uuid.New()
	rRepo.On("FindLastOpenReceptionByPVZ", pvzID).Return(nil, errors.New("not found"))

	product, err := svc.CreateProduct("обувь", pvzID)
	assert.Nil(t, product)
	assert.EqualError(t, err, "нет активной приёмки")
	rRepo.AssertExpectations(t)
}

func TestDeleteLastProduct_Success(t *testing.T) {
	pRepo := new(mocks.ProductRepo)
	rRepo := new(mocks.ReceptionRepo)
	svc := service.NewProductService(pRepo, rRepo)

	pvzID := uuid.New()
	reception := &model.Reception{ID: uuid.New()}
	product := &model.Product{ID: uuid.New()}

	rRepo.On("FindLastOpenReceptionByPVZ", pvzID).Return(reception, nil)
	pRepo.On("GetLastAddedProduct", reception.ID).Return(product, nil)
	pRepo.On("DeleteProductByID", product.ID).Return(nil)

	err := svc.DeleteLastProduct(pvzID, "employee")
	assert.NoError(t, err)
	pRepo.AssertExpectations(t)
	rRepo.AssertExpectations(t)
}

func TestDeleteLastProduct_NotEmployee(t *testing.T) {
	svc := service.NewProductService(nil, nil)
	err := svc.DeleteLastProduct(uuid.New(), "moderator")
	assert.EqualError(t, err, "доступ запрещён")
}

func TestDeleteLastProduct_NoReception(t *testing.T) {
	rRepo := new(mocks.ReceptionRepo)
	svc := service.NewProductService(nil, rRepo)

	pvzID := uuid.New()
	rRepo.On("FindLastOpenReceptionByPVZ", pvzID).Return(nil, errors.New("not found"))

	err := svc.DeleteLastProduct(pvzID, "employee")
	assert.EqualError(t, err, "не удалось найти последнюю открытую приемку")
	rRepo.AssertExpectations(t)
}

func TestDeleteLastProduct_NoProduct(t *testing.T) {
	pRepo := new(mocks.ProductRepo)
	rRepo := new(mocks.ReceptionRepo)
	svc := service.NewProductService(pRepo, rRepo)

	pvzID := uuid.New()
	reception := &model.Reception{ID: uuid.New()}

	rRepo.On("FindLastOpenReceptionByPVZ", pvzID).Return(reception, nil)
	pRepo.On("GetLastAddedProduct", reception.ID).Return(nil, errors.New("not found"))

	err := svc.DeleteLastProduct(pvzID, "employee")
	assert.EqualError(t, err, "не удалось найти последний добавленный товар")
	pRepo.AssertExpectations(t)
	rRepo.AssertExpectations(t)
}

func TestDeleteLastProduct_DeleteError(t *testing.T) {
	pRepo := new(mocks.ProductRepo)
	rRepo := new(mocks.ReceptionRepo)
	svc := service.NewProductService(pRepo, rRepo)

	pvzID := uuid.New()
	reception := &model.Reception{ID: uuid.New()}
	product := &model.Product{ID: uuid.New()}

	rRepo.On("FindLastOpenReceptionByPVZ", pvzID).Return(reception, nil)
	pRepo.On("GetLastAddedProduct", reception.ID).Return(product, nil)
	pRepo.On("DeleteProductByID", product.ID).Return(errors.New("fail"))

	err := svc.DeleteLastProduct(pvzID, "employee")
	assert.EqualError(t, err, "не удалось удалить товар")
	pRepo.AssertExpectations(t)
	rRepo.AssertExpectations(t)
}
