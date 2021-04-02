package data

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func SignInRequest(requestBody []byte) []byte {
	signinPath := "http://localhost:8080/auth/realms/ubivius/protocol/openid-connect/token"

	username := ExtractValue(string(requestBody), "username")
	password := ExtractValue(string(requestBody), "password")
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
	return body
}

func SignUpRequest(requestBody []byte) (string, string) {
	newUserPath := "http://localhost:8080/auth/admin/realms/ubivius/users"

	playerData := ExtractPlayerData(requestBody)
	jsonValues := AddValueToList(playerData, "enabled", "true")

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

	return resp.Status, admin_token
}

func AddNewUser(playerId string, requestBody []byte) {
	addUserPath := "http://localhost:9091/users"

	jsonValues := ExtractPlayerData(requestBody)

	newPlayer := AddValueToList(jsonValues, "id", playerId)

	req, err := http.NewRequest("POST", addUserPath, bytes.NewBuffer(newPlayer))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func SetUserPassword(playerId string, password string, admin_token string) {
	userSetPasswordPath := "http://localhost:8080/auth/admin/realms/ubivius/users/" + playerId + "/reset-password"

	values := map[string]string{"type": "password", "temporary": "false", "value": password}
	jsonValues, _ := json.Marshal(values)

	req, err := http.NewRequest("PUT", userSetPasswordPath, bytes.NewBuffer(jsonValues))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+admin_token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
func GetUser(userId string) []byte {
	userPath := "http://localhost:9091/users/" + userId

	req, err := http.NewRequest("GET", userPath, nil)
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
	return body
}

func GetUserId(username string, admin_token string) string {
	userIdPath := "http://localhost:8080/auth/admin/realms/ubivius/users?username=" + username

	req, err := http.NewRequest("GET", userIdPath, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+admin_token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	playerId := ExtractValue(string(body), "id")

	return playerId
}

func GetAdminAccessToken() string {
	adminAccessPath := "http://localhost:8080/auth/realms/ubivius/protocol/openid-connect/token"

	data := url.Values{}
	data.Set("client_id", "ubivius-client")
	data.Set("grant_type", "client_credentials")
	data.Set("client_secret", "7d109d2b-524f-4351-bfda-44ecad030eef")

	req, err := http.NewRequest("POST", adminAccessPath, strings.NewReader(data.Encode()))
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
	admin_token := ExtractValue(string(body), "access_token")
	return admin_token
}

// extracts the value for a key from a JSON-formatted string
// body - the JSON-response as a string. Usually retrieved via the request body
// key - the key for which the value should be extracted
// returns - the value for the given key
func ExtractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}

func ExtractPlayerData(requestBody []byte) []byte {
	username := ExtractValue(string(requestBody), "username")
	firstName := ExtractValue(string(requestBody), "firstName")
	lastName := ExtractValue(string(requestBody), "lastName")
	email := ExtractValue(string(requestBody), "email")

	values := map[string]string{"firstName": firstName, "lastName": lastName, "email": email, "username": username}
	jsonValues, _ := json.Marshal(values)
	return jsonValues
}

//Extract a specific claim for the jwt token
func ExtractClaims(tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, nil)
	if err != nil {
		log.Println("Error while getting claims: ", err)
	}
	return claims
}

func AddValueToList(values []byte, newKey string, newValue string) []byte {
	jsonBody := map[string]string{}
	err := json.Unmarshal(values, &jsonBody)
	if err != nil {
		panic(err)
	}
	jsonBody[newKey] = newValue
	newList, _ := json.Marshal(jsonBody)
	return newList
}
