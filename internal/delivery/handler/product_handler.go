package handler

import (
	"encoding/json"
	"net/http"
	"pvz/internal/delivery/middleware"
	"pvz/internal/domain/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

type createProductRequest struct {
	Type  string `json:"type"`
	PVZID string `json:"pvzId"`
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "employee" {
		http.Error(w, `{"message":"доступ запрещен"}`, http.StatusForbidden)
		return
	}

	var req createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"неверный запрос"}`, http.StatusBadRequest)
		return
	}

	product, err := h.service.CreateProduct(req.Type, uuid.MustParse(req.PVZID))
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	role, err := middleware.GetUserRole(r.Context())
	if err != nil || role != "employee" {
		http.Error(w, `{"message":"доступ запрещен"}`, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	pvzID := uuid.MustParse(vars["pvzId"])

	err = h.service.DeleteLastProduct(pvzID, role)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
