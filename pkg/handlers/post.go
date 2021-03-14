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

	claims := extractClaims(access_token)
	// do something with decoded claims
	/*for key, val := range claims {
		log.Print(val)
		log.Print(key)
	}*/
	log.Print(claims["sub"])
	userId := claims["sub"]

	player := map[string]string{"userId": string(userId), "access_token": access_token}
	playerJson, _ := json.Marshal(player)

	err = json.NewEncoder(responseWriter).Encode(playerJson)
	if err != nil {
		http.Error(responseWriter, "Unable to complete SignIn", http.StatusInternalServerError)
	}
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
	log.Println(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	access_token := extractValue(string(body), "access_token")
	log.Println("response Body:", string(body))

	player := map[string]string{"username": username, "access_token": access_token}
	playerJson, _ := json.Marshal(player)

	err = json.NewEncoder(responseWriter).Encode(playerJson)
	if err != nil {
		http.Error(responseWriter, "Unable to complete SignIn", http.StatusInternalServerError)
	}
}
