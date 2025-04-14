package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getAuthToken(t *testing.T, serverURL string, role string) string {
	loginData := map[string]string{
		"role": role,
	}
	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Failed to marshal login data: %v", err)
	}

	resp, err := http.Post(serverURL+"/dummyLogin", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatalf("Failed to send login request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}
	token, ok := result["token"].(string)
	if !ok {
		t.Fatalf("Token not found in login response")
	}

	return token
}

func createPVZ(t *testing.T, serverURL, token, city string) string {
	reqBody := map[string]string{
		"city": city,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal city: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, serverURL+"/pvz", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 201, got %d. Response body: %s", resp.StatusCode, body)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	idValue, ok := result["ID"]
	if !ok {
		t.Fatalf("Response JSON missing 'id' field: %v", result)
	}

	pvzID, ok := idValue.(string)
	if !ok {
		t.Fatalf("Field 'id' is not a string, got: %T", idValue)
	}

	return pvzID
}

func createReception(t *testing.T, serverURL, token, pvzID string) string {
	reqBody := map[string]string{
		"pvzId": pvzID,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal reception: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, serverURL+"/receptions", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	receptionID := result["ID"].(string)
	return receptionID
}

func addProduct(t *testing.T, serverURL, token, PVZID, productType string) {
	reqBody := map[string]string{
		"PVZId": PVZID,
		"type":  productType,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal product: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, serverURL+"/products", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func closeReception(t *testing.T, serverURL, token, pvzID string) {
	req, err := http.NewRequest(http.MethodPost, serverURL+"/pvz/"+pvzID+"/close_last_reception", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFullFlow(t *testing.T) {
	serverURL := "http://localhost:8080"
	token_emp := getAuthToken(t, serverURL, "employee")
	token_mod := getAuthToken(t, serverURL, "moderator")

	pvzID := createPVZ(t, serverURL, token_mod, "Москва")
	_ = createReception(t, serverURL, token_emp, pvzID)
	for i := 1; i <= 50; i++ {
		productType := "обувь"
		addProduct(t, serverURL, token_emp, pvzID, productType)
	}
	closeReception(t, serverURL, token_emp, pvzID)
}
