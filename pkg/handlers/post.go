package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/Ubivius/microservice-authentication/pkg/data"
)

// SignIn will fetch the acces token from Keycloak and return user data
func (authHandler *AuthHandler) SignIn(responseWriter http.ResponseWriter, request *http.Request) {

	log.Info("SignIn request")
	requestBody, _ := ioutil.ReadAll(request.Body)

	//SignIn request
	signInBody := data.SignInRequest(requestBody)

	//Get access token
	access_token := data.ExtractValue(string(signInBody), "access_token")

	//Get user data
	claims := data.ExtractClaims(access_token)
	userId := claims["sub"]
	userBody := data.GetUser(userId.(string))
	player := data.AddValueToList(userBody, "accessToken", access_token)

	responseWriter.WriteHeader(http.StatusOK)
	_, err := responseWriter.Write(player)
	if err != nil {
		panic(err)
	}
}

// SignUp will register a new user in keycloak and in our user database
func (authHandler *AuthHandler) SignUp(responseWriter http.ResponseWriter, request *http.Request) {

	log.Info("SignUp request")
	requestBody, _ := ioutil.ReadAll(request.Body)

	signupStatus, admin_token := data.SignUpRequest(requestBody)

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

	//Get user Id
	playerId := data.GetUserId(data.ExtractValue(string(requestBody), "username"), admin_token)

	//Set user password
	data.SetUserPassword(playerId, data.ExtractValue(string(requestBody), "password"), admin_token)

	//Add new user
	data.AddNewUser(playerId, requestBody)
}
