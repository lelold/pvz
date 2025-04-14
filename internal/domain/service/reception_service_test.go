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

func TestStartReception(t *testing.T) {
	mockRepo := new(mocks.ReceptionRepo)
	svc := service.NewReceptionService(mockRepo)

	pvzID := uuid.New()

	t.Run("успешный старт приёмки", func(t *testing.T) {
		mockRepo.On("HasOpenReception", pvzID).Return(false, nil).Once()
		mockRepo.On("CreateReception", mock.AnythingOfType("*model.Reception")).Return(nil).Once()

		reception, err := svc.StartReception(pvzID)
		assert.NoError(t, err)
		assert.NotNil(t, reception)
		assert.Equal(t, "in_progress", reception.Status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("приёмка уже открыта", func(t *testing.T) {
		mockRepo.On("HasOpenReception", pvzID).Return(true, nil).Once()

		reception, err := svc.StartReception(pvzID)
		assert.Error(t, err)
		assert.Nil(t, reception)
		assert.Equal(t, "предыдущая приёмка ещё не закрыта", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("ошибка при проверке открытых приёмок", func(t *testing.T) {
		mockRepo.On("HasOpenReception", pvzID).Return(false, errors.New("db error")).Once()

		reception, err := svc.StartReception(pvzID)
		assert.Error(t, err)
		assert.Nil(t, reception)

		mockRepo.AssertExpectations(t)
	})
}

func TestCloseLastReception(t *testing.T) {
	mockRepo := new(mocks.ReceptionRepo)
	svc := service.NewReceptionService(mockRepo)

	pvzID := uuid.New()
	reception := &model.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		DateTime: time.Now(),
		Status:   "in_progress",
	}

	t.Run("успешное закрытие приёмки", func(t *testing.T) {
		mockRepo.On("FindLastOpenReceptionByPVZ", pvzID).Return(reception, nil).Once()
		mockRepo.On("CloseReception", reception).Return(nil).Once()

		result, err := svc.CloseLastReception(pvzID)
		assert.NoError(t, err)
		assert.Equal(t, reception.ID, result.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ошибка при поиске открытой приёмки", func(t *testing.T) {
		mockRepo.On("FindLastOpenReceptionByPVZ", pvzID).Return(nil, errors.New("not found")).Once()

		result, err := svc.CloseLastReception(pvzID)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ошибка при закрытии приёмки", func(t *testing.T) {
		mockRepo.On("FindLastOpenReceptionByPVZ", pvzID).Return(reception, nil).Once()
		mockRepo.On("CloseReception", reception).Return(errors.New("db error")).Once()

		result, err := svc.CloseLastReception(pvzID)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})
}
