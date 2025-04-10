package pvz_http

import (
	"encoding/json"
	"net/http"
	"pvz/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

var pvzStore []model.PVZ

var allowedCities = map[string]bool{
	"Москва":          true,
	"Санкт-Петербург": true,
	"Казань":          true,
}

func PVZHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/pvz", handlePVZ)
	return mux
}

func handlePVZ(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createPVZ(w, r)
	case http.MethodGet:
		getPVZList(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func createPVZ(w http.ResponseWriter, r *http.Request) {
	role, err := GetUserRole(r.Context())
	if err != nil || role != "moderator" {
		http.Error(w, `{"message":"доступ запрещен"}`, http.StatusForbidden)
		return
	}

	var req model.PVZ
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !allowedCities[req.City] {
		http.Error(w, `{"message":"неверный запрос"}`, http.StatusBadRequest)
		return
	}

	req.RegistrationDate = time.Now()
	req.ID = uuid.New()
	pvzStore = append(pvzStore, req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func getPVZList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pvzStore)
}
