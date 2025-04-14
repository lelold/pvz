package service_test

import (
	"errors"
	"pvz/internal/domain/mocks"
	"pvz/internal/domain/model"
	"pvz/internal/domain/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePVZ(t *testing.T) {
	mockRepo := new(mocks.PVZRepo)
	svc := service.NewPVZService(mockRepo)

	t.Run("успешное создание ПВЗ", func(t *testing.T) {
		city := "Москва"
		mockRepo.On("Create", mock.AnythingOfType("*model.PVZ")).Return(nil).Once()

		pvz, err := svc.CreatePVZ(city)

		assert.NoError(t, err)
		assert.NotNil(t, pvz)
		assert.Equal(t, city, pvz.City)
		mockRepo.AssertExpectations(t)
	})

	t.Run("неразрешённый город", func(t *testing.T) {
		pvz, err := svc.CreatePVZ("Новосибирск")
		assert.Error(t, err)
		assert.Nil(t, pvz)
		assert.Equal(t, "город не разрешен", err.Error())
	})
}

func TestGetFullPVZList(t *testing.T) {
	mockRepo := new(mocks.PVZRepo)
	svc := service.NewPVZService(mockRepo)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	page := 1
	limit := 10

	pvzID := uuid.New()
	pvz := model.PVZ{
		ID:               pvzID,
		City:             "Москва",
		RegistrationDate: time.Now(),
	}
	reception := model.ReceptionWithProducts{
		Reception: model.Reception{
			ID:     uuid.New(),
			PVZID:  pvzID,
			Status: "in_progress",
		},
		Products: []model.Product{
			{ID: uuid.New(), Type: "одежда"},
		},
	}

	t.Run("успешное получение списка", func(t *testing.T) {
		mockRepo.On("GetFilteredPVZs", &start, &end, page, limit).
			Return([]model.PVZ{pvz}, nil).Once()

		mockRepo.On("GetReceptionsWithProducts", pvzID).
			Return([]model.ReceptionWithProducts{reception}, nil).Once()

		result, err := svc.GetFullPVZList(&start, &end, page, limit)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, pvz.ID, result[0].PVZ.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ошибка при получении ПВЗ", func(t *testing.T) {
		mockRepo.On("GetFilteredPVZs", &start, &end, page, limit).
			Return(nil, errors.New("db error")).Once()

		result, err := svc.GetFullPVZList(&start, &end, page, limit)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
