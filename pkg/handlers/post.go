package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// SignIn will fetch the acces token from Keycloak and return user data
func (authHandler *AuthHandler) SignIn(responseWriter http.ResponseWriter, request *http.Request) {
	signinPath := "http://localhost:8080/auth/realms/ubivius/protocol/openid-connect/token"

	requestbody, _ := ioutil.ReadAll(request.Body)

	username := ExtractValue(string(requestbody), "username")
	password := ExtractValue(string(requestbody), "password")
	data := url.Values{}
	data.Set("client_id", "ubivius-client")
	data.Set("grant_type", "password")
	data.Set("client_secret", "7d109d2b-524f-4351-bfda-44ecad030eef")
	data.Set("scope", "openid")
	data.Set("username", username)
	data.Set("password", password)

	//SignIn
	req, err := http.NewRequest("POST", signinPath, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	access_token := ExtractValue(string(body), "access_token")
	//Get user data
	claims := ExtractClaims(access_token)
	userId := claims["sub"]

	userPath := "http://localhost:9091/users/" + userId.(string)

	req, err = http.NewRequest("GET", userPath, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	jsonBody := map[string]string{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		panic(err)
	}
	jsonBody["accessToken"] = access_token
	player, _ := json.Marshal(jsonBody)

	log.Println(string(player))
	responseWriter.WriteHeader(http.StatusOK)
	_, err = responseWriter.Write(player)
	if err != nil {
		panic(err)
	}
}

// SignUp will register a new user in keycloak and in our user database
func (authHandler *AuthHandler) SignUp(responseWriter http.ResponseWriter, request *http.Request) {
	newUserPath := "http://localhost:8080/auth/admin/realms/ubivius/users"

	requestbody, _ := ioutil.ReadAll(request.Body)
	username := ExtractValue(string(requestbody), "username")
	password := ExtractValue(string(requestbody), "password")
	log.Println(string(password))
	firstName := ExtractValue(string(requestbody), "firstName")
	lastName := ExtractValue(string(requestbody), "lastName")
	email := ExtractValue(string(requestbody), "email")

	values := map[string]string{"firstName": firstName, "lastName": lastName, "email": email, "username": username, "enabled": "true"}
	jsonValues, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", newUserPath, bytes.NewBuffer(jsonValues))
	if err != nil {
		panic(err)
	}
	admin_token := GetAdminAccessToken()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+admin_token)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)

	if resp.Status == "409 Conflict" {
		responseWriter.WriteHeader(http.StatusConflict)
		_, err = responseWriter.Write([]byte("409 Conflict"))
		if err != nil {
			panic(err)
		}
		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
	_, err = responseWriter.Write([]byte("201 Created"))
	if err != nil {
		panic(err)
	}

	//Get new user ID
	userIdPath := "http://localhost:8080/auth/admin/realms/ubivius/users?username=" + username

	req, err = http.NewRequest("GET", userIdPath, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+admin_token)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	playerId := ExtractValue(string(body), "id")

	//Set Password
	userSetPasswordPath := "http://localhost:8080/auth/admin/realms/ubivius/users/" + playerId + "/reset-password"

	values = map[string]string{"type": "password", "temporary": "false", "value": password}
	jsonValues, _ = json.Marshal(values)

	req, err = http.NewRequest("PUT", userSetPasswordPath, bytes.NewBuffer(jsonValues))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+admin_token)
	req.Header.Add("Content-Type", "application/json")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//Post new user
	addUserPath := "http://localhost:9091/users"

	jsonBody := map[string]string{}
	err = json.Unmarshal(jsonValues, &jsonBody)
	if err != nil {
		panic(err)
	}
	jsonBody["id"] = playerId
	newPlayer, _ := json.Marshal(jsonBody)

	req, err = http.NewRequest("POST", addUserPath, bytes.NewBuffer(newPlayer))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
