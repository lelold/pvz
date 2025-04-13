package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz/internal/delivery/handler"
	"pvz/internal/delivery/middleware"
	"pvz/internal/domain/mocks"
	"pvz/internal/domain/model"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePVZ_Success(t *testing.T) {
	mockService := new(mocks.PVZService)
	h := handler.NewPVZHandler(mockService)

	reqBody := `{"city":"Москва"}`
	pvz := &model.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Москва",
	}

	mockService.On("CreatePVZ", "Москва").Return(pvz, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(reqBody))
	req = req.WithContext(middleware.SetRoleContext(req.Context(), "moderator"))
	w := httptest.NewRecorder()

	h.HandlePVZ(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.PVZ
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Москва", resp.City)

	mockService.AssertExpectations(t)
}

func TestCreatePVZ_Forbidden(t *testing.T) {
	mockService := new(mocks.PVZService)
	h := handler.NewPVZHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/pvz", strings.NewReader(`{"city":"Москва"}`))
	req = req.WithContext(middleware.SetRoleContext(req.Context(), "employee"))
	w := httptest.NewRecorder()

	h.HandlePVZ(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetPVZList_Success(t *testing.T) {
	mockService := new(mocks.PVZService)
	handler := handler.NewPVZHandler(mockService)

	pvzID := uuid.New()
	receptionID := uuid.New()
	now := time.Now()

	mockPVZList := []model.PVZFull{
		{
			PVZ: model.PVZ{
				ID:               pvzID,
				City:             "Москва",
				RegistrationDate: now,
			},
			Receptions: []model.ReceptionWithProducts{
				{
					Reception: model.Reception{
						ID:       receptionID,
						PVZID:    pvzID,
						DateTime: now,
						Status:   "in_progress",
					},
					Products: []model.Product{
						{
							ID:          uuid.New(),
							ReceptionID: receptionID,
							DateTime:    now,
							Type:        "электроника",
						},
					},
				},
			},
		},
	}

	mockService.
		On("GetFullPVZList", mock.Anything, mock.Anything, 1, 10).
		Return(mockPVZList, nil).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	req = req.WithContext(middleware.SetRoleContext(req.Context(), "employee"))
	w := httptest.NewRecorder()

	handler.HandlePVZ(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.PVZFull
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "Москва", response[0].PVZ.City)
	assert.Len(t, response[0].Receptions, 1)
	assert.Equal(t, "in_progress", response[0].Receptions[0].Reception.Status)
	assert.Len(t, response[0].Receptions[0].Products, 1)
	assert.Equal(t, "электроника", response[0].Receptions[0].Products[0].Type)

	mockService.AssertExpectations(t)
}

func TestGetPVZList_Forbidden(t *testing.T) {
	mockService := new(mocks.PVZService)
	h := handler.NewPVZHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	req = req.WithContext(middleware.SetRoleContext(req.Context(), "moderator"))
	w := httptest.NewRecorder()

	h.HandlePVZ(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}
