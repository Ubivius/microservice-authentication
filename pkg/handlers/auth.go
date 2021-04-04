package handlers

import (
	"log"
)

// AuthHandler contains the items common to all auth handler functions
type AuthHandler struct {
	logger *log.Logger
}

// NewAuthHandler returns a pointer to a AuthHandler with the logger passed as a parameter
func NewAuthHandler(logger *log.Logger) *AuthHandler {
	return &AuthHandler{logger}
}
