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

// AddProduct creates a new product from the received JSON
func (authHandler *AuthHandler) SignIn(responseWriter http.ResponseWriter, request *http.Request) {
	urlPath := "http://localhost:8080/auth/realms/ubivius/protocol/openid-connect/token"

	requestbody, _ := ioutil.ReadAll(request.Body)

	username := extractValue(string(requestbody), "username")
	password := extractValue(string(requestbody), "password")

	data := url.Values{}
	data.Set("client_id", "ubivius-client")
	data.Set("grant_type", "password")
	data.Set("client_secret", "7d109d2b-524f-4351-bfda-44ecad030eef")
	data.Set("scope", "openid")
	data.Set("username", username)
	data.Set("password", password)

	//SignIn
	req, err := http.NewRequest("POST", urlPath, strings.NewReader(data.Encode()))
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
	access_token := extractValue(string(body), "access_token")
	log.Println(string(body))

	claims := extractClaims(access_token)
	userId := claims["sub"]
	log.Println(userId)

	//Get user data
	urlPath = "http://localhost:9091/users/" + "1"

	req, err = http.NewRequest("GET", urlPath, nil)
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
	jsonBody["accessToken"] = access_token
	player, _ := json.Marshal(jsonBody)

	log.Println(string(player))
	responseWriter.Write(player)
}

func (authHandler *AuthHandler) SignUp(responseWriter http.ResponseWriter, request *http.Request) {
	urlPath := "http://localhost:8080/auth/admin/realms/ubivius/users"

	requestbody, _ := ioutil.ReadAll(request.Body)
	username := extractValue(string(requestbody), "username")
	firstName := extractValue(string(requestbody), "firstName")
	lastName := extractValue(string(requestbody), "lastName")
	email := extractValue(string(requestbody), "email")

	values := map[string]string{"firstName": firstName, "lastName": lastName, "email": email, "username": username, "enabled": "true"}
	jsonValues, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", urlPath, bytes.NewBuffer(jsonValues))
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
		responseWriter.Write([]byte("409 Conflict"))
		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
	responseWriter.Write([]byte("201 Created"))

	//Get new user ID
	urlPath = "http://localhost:8080/auth/admin/realms/ubivius/users?username=" + username

	req, err = http.NewRequest("GET", urlPath, nil)
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
	playerId := extractValue(string(body), "id")

	//Post new user
	urlPath = "http://localhost:9091/users"

	jsonBody := map[string]string{}
	err = json.Unmarshal(jsonValues, &jsonBody)
	jsonBody["id"] = playerId
	newPlayer, _ := json.Marshal(jsonBody)

	req, err = http.NewRequest("POST", urlPath, bytes.NewBuffer(newPlayer))
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

	body, _ = ioutil.ReadAll(resp.Body)
}
