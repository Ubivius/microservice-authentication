package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// AddProduct creates a new product from the received JSON
func (authHandler *AuthHandler) SignIn(responseWriter http.ResponseWriter, request *http.Request) {
	urlPath := "http://localhost:8080/auth/realms/ubivius/protocol/openid-connect/token"

	requestbody, _ := ioutil.ReadAll(request.Body)
	log.Println(string(requestbody))
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

	log.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	access_token := extractValue(string(body), "access_token")
	log.Println("access_token: ", access_token)

	player := map[string]string{"username": username, "access_token": access_token}
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

func GetAdminAccessToken() string {
	urlPath := "http://localhost:8080/auth/realms/ubivius/protocol/openid-connect/token"

	data := url.Values{}
	data.Set("client_id", "ubivius-client")
	data.Set("grant_type", "client_credentials")
	data.Set("client_secret", "7d109d2b-524f-4351-bfda-44ecad030eef")

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

	log.Println("Admintoken response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	admin_token := extractValue(string(body), "access_token")
	return admin_token
}

// extracts the value for a key from a JSON-formatted string
// body - the JSON-response as a string. Usually retrieved via the request body
// key - the key for which the value should be extracted
// returns - the value for the given key
func extractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}
