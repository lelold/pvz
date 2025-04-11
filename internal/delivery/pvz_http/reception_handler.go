package pvz_http

import (
	"encoding/json"
	"errors"
	"net/http"
	"pvz/internal/domain/model"
	"pvz/internal/domain/service"
)

type ReceptionHandler struct {
	Service *service.ReceptionService
}

func NewReceptionHandler(s *service.ReceptionService) *ReceptionHandler {
	return &ReceptionHandler{Service: s}
}

func (h *ReceptionHandler) StartReception(w http.ResponseWriter, r *http.Request) {
	role, err := GetUserRole(r.Context())
	if err != nil || role != "employee" {
		http.Error(w, `{"message":"доступ запрещен"}`, http.StatusForbidden)
		return
	}

	var req model.Reception
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"неверный запрос"}`, http.StatusBadRequest)
		return
	}

	reception, err := h.Service.StartReception(req.PVZID.String())
	if err != nil {
		if err == errors.New("предыдущая приёмка ещё не закрыта") {
			http.Error(w, `{"message":"предыдущая приёмка ещё не закрыта"}`, http.StatusConflict)
		} else {
			http.Error(w, `{"message":"ошибка сервиса"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reception)
}
