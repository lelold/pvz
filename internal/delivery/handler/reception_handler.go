package handler

import (
	"encoding/json"
	"net/http"
	"pvz/internal/delivery/middleware"
	"pvz/internal/domain/model"
	"pvz/internal/domain/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReceptionHandler struct {
	Service service.ReceptionService
}

func NewReceptionHandler(s service.ReceptionService) *ReceptionHandler {
	return &ReceptionHandler{Service: s}
}

func (h *ReceptionHandler) StartReception(w http.ResponseWriter, r *http.Request) {
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "employee" {
		http.Error(w, `{"message":"доступ запрещен"}`, http.StatusForbidden)
		return
	}

	var req model.Reception
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"неверный запрос"}`, http.StatusBadRequest)
		return
	}

	reception, err := h.Service.StartReception(req.PVZID)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reception)
}

func (h *ReceptionHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "employee" {
		http.Error(w, `{"message":"доступ запрещен"}`, http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	pvzID := vars["pvzId"]
	if pvzID == "" {
		http.Error(w, `{"message":"неверный запрос: отсутствует pvzId"}`, http.StatusBadRequest)
		return
	}

	reception, err := h.Service.CloseLastReception(uuid.MustParse(pvzID))
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reception)
}
