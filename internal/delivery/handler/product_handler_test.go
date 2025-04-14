package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pvz/internal/delivery/middleware"
	"pvz/internal/domain/mocks"
	"pvz/internal/domain/model"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	mockService := new(mocks.ProductService)
	handler := NewProductHandler(mockService)

	reqBody := createProductRequest{
		Type:  "электроника",
		PVZID: uuid.New().String(),
	}
	reqBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(reqBytes))
	w := httptest.NewRecorder()

	req = req.WithContext(middleware.SetRoleContext(req.Context(), "employee"))

	mockService.On("CreateProduct", reqBody.Type, uuid.MustParse(reqBody.PVZID)).Return(&model.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        reqBody.Type,
		ReceptionID: uuid.New(),
	}, nil).Once()

	handler.CreateProduct(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.Product
	err = json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.ID)
	assert.Equal(t, reqBody.Type, resp.Type)
	assert.NotNil(t, resp.ReceptionID)
	mockService.AssertExpectations(t)
}

func TestDeleteLastProduct(t *testing.T) {
	mockService := new(mocks.ProductService)
	handler := NewProductHandler(mockService)

	pvzID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID+"/delete_last_product", nil)
	w := httptest.NewRecorder()

	req = req.WithContext(middleware.SetRoleContext(req.Context(), "employee"))
	req = mux.SetURLVars(req, map[string]string{"pvzId": pvzID})

	mockService.On("DeleteLastProduct", uuid.MustParse(pvzID), "employee").Return(nil).Once()

	handler.DeleteLastProduct(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateProductForbiddenRole(t *testing.T) {
	mockService := new(mocks.ProductService)
	handler := NewProductHandler(mockService)

	reqBody := createProductRequest{
		Type:  "st",
		PVZID: "sid",
	}
	reqBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(reqBytes))
	w := httptest.NewRecorder()

	mockRole := "moderator"
	req = req.WithContext(middleware.SetRoleContext(req.Context(), mockRole))

	handler.CreateProduct(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var resp map[string]string
	err = json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "доступ запрещен", resp["message"])
}
