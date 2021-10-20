package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Ubivius/microservice-authentication/pkg/data"
)

// SignIn will fetch the acces token from Keycloak and return user data
func (authHandler *AuthHandler) SignIn(responseWriter http.ResponseWriter, request *http.Request) {

	log.Info("SignIn request")

	//Extract credentials
	var credentials data.Credentials
	errorCreds := json.NewDecoder(request.Body).Decode(&credentials)
    if errorCreds != nil {
        http.Error(responseWriter, errorCreds.Error(), http.StatusBadRequest)
        return
    }

	//SignIn request
	player := data.SignInRequest(credentials)

	responseWriter.WriteHeader(http.StatusOK)
	_, err := responseWriter.Write(player)
	if err != nil {
		panic(err)
	}
}

// SignUp will register a new user in keycloak and in our user database
func (authHandler *AuthHandler) SignUp(responseWriter http.ResponseWriter, request *http.Request) {

	log.Info("SignUp request")

	//Extract user info
	var user data.User
	errorUser := json.NewDecoder(request.Body).Decode(&user)
    if errorUser != nil {
        http.Error(responseWriter, errorUser.Error(), http.StatusBadRequest)
        return
    }

	//SignUp request
	admin_token := data.GetAdminAccessToken()
	signupStatus := data.SignUpRequest(user, admin_token)

	log.Info("SignUp Response:", "status", signupStatus)

	if signupStatus == "409 Conflict" {
		responseWriter.WriteHeader(http.StatusConflict)
		_, err := responseWriter.Write([]byte("409 Conflict"))
		if err != nil {
			panic(err)
		}
		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
	_, err := responseWriter.Write([]byte("201 Created"))
	if err != nil {
		panic(err)
	}
}
