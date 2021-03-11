package handlers

import (
	"log"
)

// KeyProduct is a key used for the Product object inside context
type KeyProduct struct{}

// AuthHandler contains the items common to all product handler functions
type AuthHandler struct {
	logger *log.Logger
}

// NewAuthHandler returns a pointer to a AuthHandler with the logger passed as a parameter
func NewAuthHandler(logger *log.Logger) *AuthHandler {
	return &AuthHandler{logger}
}
