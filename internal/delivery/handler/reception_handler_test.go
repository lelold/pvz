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
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestStartReception_Success(t *testing.T) {
	mockService := new(mocks.ReceptionService)
	handler := handler.NewReceptionHandler(mockService)

	pvzID := uuid.New()
	expectedReception := &model.Reception{ID: uuid.New(), PVZID: pvzID, Status: "in_progress"}
	mockService.On("StartReception", pvzID).Return(expectedReception, nil).Once()

	body, _ := json.Marshal(model.Reception{PVZID: pvzID})
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(body))
	req = req.WithContext(middleware.SetRoleContext(req.Context(), "employee"))

	w := httptest.NewRecorder()
	handler.StartReception(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var actual model.Reception
	json.NewDecoder(w.Body).Decode(&actual)
	assert.Equal(t, expectedReception.ID, actual.ID)

	mockService.AssertExpectations(t)
}

func TestCloseLastReception_Success(t *testing.T) {
	mockService := new(mocks.ReceptionService)
	handler := handler.NewReceptionHandler(mockService)

	pvzID := uuid.New().String()
	expected := &model.Reception{ID: uuid.New(), PVZID: uuid.MustParse(pvzID), Status: "close"}

	mockService.On("CloseLastReception", uuid.MustParse(pvzID)).Return(expected, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID+"/close_last_reception", nil)
	req = req.WithContext(middleware.SetRoleContext(req.Context(), "employee"))
	req = mux.SetURLVars(req, map[string]string{
		"pvzId": pvzID,
	})

	w := httptest.NewRecorder()
	handler.CloseLastReception(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var actual model.Reception
	err := json.NewDecoder(w.Body).Decode(&actual)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, actual.ID)

	mockService.AssertExpectations(t)
}
