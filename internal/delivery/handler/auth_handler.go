package handler

import (
	"encoding/json"
	"net/http"

	"pvz/internal/delivery/middleware"
	"pvz/internal/domain/model"
	"pvz/internal/domain/service"

	"github.com/google/uuid"
)

type AuthHandler struct {
	UserService service.UserService
}

func DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"неверный запрос"}`, http.StatusBadRequest)
		return
	}

	if req.Role != "employee" && req.Role != "moderator" {
		http.Error(w, `{"message":"неверная роль"}`, http.StatusBadRequest)
		return
	}

	tokenString, err := middleware.GenerateToken(uuid.NewString(), req.Role)
	if err != nil {
		http.Error(w, `{"message":"ошибка с генерацией токена"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.TokenResponse{Token: tokenString})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"неверный запрос"}`, http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Register(req.Email, req.Password, req.Role)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"id":    user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"неверный запрос"}`, http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"message":"неавторизован"}`, http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		http.Error(w, `{"message":"ошибка с генерацией токена"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
