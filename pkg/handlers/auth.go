package handlers

// AuthHandler contains the items common to all auth handler functions
type AuthHandler struct {}

// NewAuthHandler returns a pointer to a AuthHandler with the logger passed as a parameter
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}
