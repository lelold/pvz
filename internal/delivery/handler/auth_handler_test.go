package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pvz/internal/delivery/handler"
	"pvz/internal/domain/mocks"
	"pvz/internal/domain/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDummyLoginHandler(t *testing.T) {
	reqBody := model.DummyLoginRequest{
		Role: "employee",
	}

	reqBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(reqBytes))
	w := httptest.NewRecorder()

	middleware := handler.DummyLoginHandler
	middleware(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err = json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp["token"])
}

func TestRegister(t *testing.T) {
	mockUserService := new(mocks.UserService)
	handler := handler.AuthHandler{
		UserService: mockUserService,
	}

	reqBody := model.RegisterRequest{
		Email:    "test@example.com",
		Password: "123",
		Role:     "employee",
	}

	reqBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBytes))
	w := httptest.NewRecorder()

	mockUserService.On("Register", reqBody.Email, reqBody.Password, reqBody.Role).
		Return(model.User{
			ID:    uuid.New(),
			Email: reqBody.Email,
			Role:  reqBody.Role,
		}, nil)

	handler.Register(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err = json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.Email, resp["email"])
	assert.Equal(t, reqBody.Role, resp["role"])

	mockUserService.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockUserService := new(mocks.UserService)
	handler := handler.AuthHandler{
		UserService: mockUserService,
	}

	reqBody := model.LoginRequest{
		Email:    "test@example.com",
		Password: "123",
	}

	reqBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBytes))
	w := httptest.NewRecorder()

	mockUserService.On("Login", reqBody.Email, reqBody.Password).
		Return(model.User{
			ID:    uuid.New(),
			Email: reqBody.Email,
			Role:  "employee",
		}, nil)

	handler.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err = json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp["token"])

	mockUserService.AssertExpectations(t)
}
