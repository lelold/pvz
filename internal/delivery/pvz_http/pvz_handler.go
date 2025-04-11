package pvz_http

import (
	"encoding/json"
	"net/http"
	"pvz/internal/domain/service"
)

type PVZHandler struct {
	service *service.PVZService
}

func NewPVZHandler(service *service.PVZService) *PVZHandler {
	return &PVZHandler{service: service}
}

type createPVZRequest struct {
	City string `json:"city"`
}

func (h *PVZHandler) HandlePVZ(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createPVZ(w, r)
	case http.MethodGet:
		h.getAll(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *PVZHandler) createPVZ(w http.ResponseWriter, r *http.Request) {
	role, err := GetUserRole(r.Context())
	if err != nil || role != "moderator" {
		http.Error(w, `{"message":"доступ запрещен"}`, http.StatusForbidden)
		return
	}

	var req createPVZRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"неверный запрос"}`, http.StatusBadRequest)
		return
	}

	pvz, err := h.service.CreatePVZ(req.City)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pvz)
}

func (h *PVZHandler) getAll(w http.ResponseWriter, r *http.Request) {
	pvzList, err := h.service.GetAllPVZ()
	if err != nil {
		http.Error(w, `{"message":"ошибка получения ПВЗ"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pvzList)
}
