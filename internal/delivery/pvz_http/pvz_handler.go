package pvz_http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type Reception struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    string    `json:"pvzId"`
	Status   string    `json:"status"`
}

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}

var pvzStore []PVZ

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

	var req PVZ
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.City == "" {
		http.Error(w, `{"message":"неверный запрос"}`, http.StatusBadRequest)
		return
	}

	req.RegistrationDate = time.Now()
	req.ID = uuid.New().String()
	pvzStore = append(pvzStore, req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func getPVZList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pvzStore)
}
