package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pvz/internal/delivery/handler"
	"pvz/internal/domain/mocks"
	"pvz/internal/domain/model"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func withRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, "role", role)
}

func TestCreateProduct_Success(t *testing.T) {
	mockService := new(mocks.ProductService)
	h := handler.NewProductHandler(mockService)

	pvzID := uuid.New()
	reqBody := map[string]string{"type": "electronics", "pvzId": pvzID.String()}
	body, _ := json.Marshal(reqBody)

	r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	r = r.WithContext(withRole(r.Context(), "employee"))
	w := httptest.NewRecorder()

	expectedProduct := model.Product{ID: uuid.New()}
	mockService.On("CreateProduct", "electronics", pvzID).Return(expectedProduct, nil)

	h.CreateProduct(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateProduct_Forbidden(t *testing.T) {
	h := handler.NewProductHandler(nil)

	r := httptest.NewRequest(http.MethodPost, "/products", nil)
	r = r.WithContext(withRole(r.Context(), "client"))
	w := httptest.NewRecorder()

	h.CreateProduct(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreateProduct_InvalidBody(t *testing.T) {
	h := handler.NewProductHandler(nil)

	r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer([]byte("{invalid")))
	r = r.WithContext(withRole(r.Context(), "employee"))
	w := httptest.NewRecorder()

	h.CreateProduct(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProduct_Fail(t *testing.T) {
	mockService := new(mocks.ProductService)
	h := handler.NewProductHandler(mockService)

	pvzID := uuid.New()
	reqBody := map[string]string{"type": "electronics", "pvzId": pvzID.String()}
	body, _ := json.Marshal(reqBody)

	r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	r = r.WithContext(withRole(r.Context(), "employee"))
	w := httptest.NewRecorder()

	mockService.On("CreateProduct", "electronics", pvzID).Return(model.Product{}, errors.New("some error"))

	h.CreateProduct(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteLastProduct_Success(t *testing.T) {
	mockService := new(mocks.ProductService)
	h := handler.NewProductHandler(mockService)

	pvzID := uuid.New()
	r := httptest.NewRequest(http.MethodDelete, "/products/"+pvzID.String(), nil)
	r = mux.SetURLVars(r, map[string]string{"pvzId": pvzID.String()})
	r = r.WithContext(withRole(r.Context(), "employee"))
	w := httptest.NewRecorder()

	mockService.On("DeleteLastProduct", pvzID, "employee").Return(nil)

	h.DeleteLastProduct(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteLastProduct_Forbidden(t *testing.T) {
	h := handler.NewProductHandler(nil)

	r := httptest.NewRequest(http.MethodDelete, "/products/someid", nil)
	r = r.WithContext(withRole(r.Context(), "client"))
	w := httptest.NewRecorder()

	h.DeleteLastProduct(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteLastProduct_Fail(t *testing.T) {
	mockService := new(mocks.ProductService)
	h := handler.NewProductHandler(mockService)

	pvzID := uuid.New()
	r := httptest.NewRequest(http.MethodDelete, "/products/"+pvzID.String(), nil)
	r = mux.SetURLVars(r, map[string]string{"pvzId": pvzID.String()})
	r = r.WithContext(withRole(r.Context(), "employee"))
	w := httptest.NewRecorder()

	mockService.On("DeleteLastProduct", pvzID, "employee").Return(errors.New("fail"))

	h.DeleteLastProduct(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}
