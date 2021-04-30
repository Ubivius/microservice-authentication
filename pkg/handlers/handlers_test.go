package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignInIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Test skipped during unit tests")
	}

	bodyReader := strings.NewReader(`{"username":"sickboy","password":"ubi123"}`)

	request := httptest.NewRequest(http.MethodPost, "/signin", bodyReader)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()

	productHandler := NewAuthHandler()
	productHandler.SignIn(response, request)

	if response.Code != 200 {
		t.Errorf("Expected status code 200 but got : %d", response.Code)
	}
	if !strings.Contains(response.Body.String(), "\"id\":2") {
		t.Error("Missing elements from expected results")
	}
}
