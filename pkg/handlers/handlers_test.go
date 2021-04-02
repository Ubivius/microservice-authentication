package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Move to util package in Sprint 9, should be a testing specific logger
func NewTestLogger() *log.Logger {
	return log.New(os.Stdout, "Tests", log.LstdFlags)
}

func TestSignIn(t *testing.T) {
	bodyReader := strings.NewReader(`{"username":"sickboy","password":"ubi123"}`)

	request := httptest.NewRequest(http.MethodPost, "/signin", bodyReader)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()

	productHandler := NewAuthHandler(NewTestLogger())
	productHandler.SignIn(response, request)

	if response.Code != 200 {
		t.Errorf("Expected status code 200 but got : %d", response.Code)
	}
	if !strings.Contains(response.Body.String(), "\"id\":2") {
		t.Error("Missing elements from expected results")
	}
}
