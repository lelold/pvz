package handler

import (
	"encoding/json"
	"net/http"
	"pvz/internal/delivery/middleware"
	"pvz/internal/domain/service"
	"strconv"
	"time"
)

type PVZHandler struct {
	service service.PVZService
}

func NewPVZHandler(service service.PVZService) *PVZHandler {
	return &PVZHandler{service: service}
}

type createPVZRequest struct {
	City string `json:"city"`
}

func (h *PVZHandler) HandlePVZ(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		role, err := middleware.GetUserRole(r.Context())
		if err != nil || role != "moderator" {
			http.Error(w, `{"message":"доступ запрещен"}`, http.StatusForbidden)
			return
		}
		h.createPVZ(w, r)

	case http.MethodGet:
		role, err := middleware.GetUserRole(r.Context())
		if err != nil || role != "employee" {
			http.Error(w, `{"message":"доступ запрещен"}`, http.StatusForbidden)
			return
		}
		h.getPVZList(w, r)

	default:
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}

func (h *PVZHandler) createPVZ(w http.ResponseWriter, r *http.Request) {
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

func (h *PVZHandler) getPVZList(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("startDate")
	endStr := r.URL.Query().Get("endDate")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 30 {
		limit = 10
	}

	var startTime, endTime *time.Time
	if startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = &t
		}
	}
	if endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			endTime = &t
		}
	}

	pvzList, err := h.service.GetFullPVZList(startTime, endTime, page, limit)
	if err != nil {
		http.Error(w, `{"message":"Не удалось получить список ПВЗ"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pvzList)
}
