package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// AuthHandler contains the items common to all auth handler functions
type AuthHandler struct {
	logger *log.Logger
}

// NewAuthHandler returns a pointer to a AuthHandler with the logger passed as a parameter
func NewAuthHandler(logger *log.Logger) *AuthHandler {
	return &AuthHandler{logger}
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

//Not used yet, but could be usefull
func ExtractClaims(tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, nil)
	if err != nil {
		log.Println("Error while getting claims: ", err)
	}
	return claims
}
