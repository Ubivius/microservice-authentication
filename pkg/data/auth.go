package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

const KeycloakPath = "http://keycloak"
const MicroserviceUserPath = "http://microservice-user:9090"

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type User struct {
	ID           string `json:"id" bson:"_id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	DateOfBirth  string `json:"dateofbirth"`
	Gender       string `json:"gender"`
	Address      string `json:"address"`
	Bio          string `json:"bio"`
	Achievements string `json:"achievements"`
}

type KeycloakUser struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Enabled      bool   `json:"enabled"`	
}

/* Claims Map 
	AuthContextClassReference  acr
	AuthorizedParty            azp
	Email                      email
	EmailVerified              email_verified
	ExpirationTime             exp
	LastName                   family_name
	FirstName                  given_name
	IssuedAt                   iat
	Issuer                     iss
	JwtID                      jti
	FullName                   name
	Username                   preferred_username
	Scope                      scope
	SessionState               session_state
	SessionID                  sid
	Subject                    sub
	Type                       typ
*/

// ErrorEnvVar : Environment variable error
var ErrorEnvVar = fmt.Errorf("missing environment variable")

func SignInRequest(credentials Credentials) []byte {
	tokenResponse := GetAccessToken(credentials)

	//Get user data
	claims := ExtractClaims(tokenResponse.AccessToken)
	userId := claims["sub"]
	userBody := GetUser(userId.(string))
	player := AddValueToList(userBody, "accessToken", tokenResponse.AccessToken)

	return player
}

func SignUpRequest(user User, admin_token string) string {
	newUserPath := KeycloakPath + "/auth/admin/realms/master/users"

	userJSON, errorUser := json.Marshal(UserToKeycloakUser(user))
    if errorUser != nil {
        log.Error(errorUser, "Error serializing user")
    }

	req, err := http.NewRequest("POST", newUserPath, bytes.NewBuffer(userJSON))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + admin_token)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	SetUserPassword(location, user.Password, admin_token)
	
	locationSplit := strings.Split(location,"/")
	user.ID = locationSplit[len(locationSplit)-1]
	AddNewUser(user)

	return resp.Status
}

func AddNewUser(user User) {
	addUserPath := MicroserviceUserPath + "/users"

	userJSON, err := json.Marshal(user)
    if err != nil {
        log.Error(err, "Error serializing user")
    }

	req, err := http.NewRequest("POST", addUserPath, bytes.NewBuffer(userJSON))
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

func SetUserPassword(path string, password string, admin_token string) {
	userSetPasswordPath := path + "/reset-password"

	values := map[string]string{"type": "password", "temporary": "false", "value": password}
	jsonValues, _ := json.Marshal(values)

	req, err := http.NewRequest("PUT", userSetPasswordPath, bytes.NewBuffer(jsonValues))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer " + admin_token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func GetUser(userId string) []byte {
	userPath := MicroserviceUserPath + "/users/" + userId

	req, err := http.NewRequest("GET", userPath, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func GetAccessToken(credentials Credentials) TokenResponse {
	signinPath := KeycloakPath + "/auth/realms/master/protocol/openid-connect/token"

	data := url.Values{}
	data.Set("client_id", "admin-cli")
	data.Set("grant_type", "password")
	data.Set("username", credentials.Username)
	data.Set("password", credentials.Password)

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

	var tokenResponse TokenResponse
	errorToken := json.NewDecoder(resp.Body).Decode(&tokenResponse)
    if errorToken != nil {
        panic(errorToken)
    }
	defer resp.Body.Close()

	return tokenResponse
}

func GetAdminAccessToken() string {
	username := os.Getenv("KEYCLOAK_ADMIN_USER")
	password := os.Getenv("KEYCLOAK_ADMIN_PASSWORD")

	if username == "" || password == "" {
		log.Error(ErrorEnvVar, "Some environment variables are not available for the Keycloak connection. KEYCLOAK_ADMIN_USER, KEYCLOAK_ADMIN_PASSWORD")
		os.Exit(1)
	}

	credentials := Credentials{
		Username: username, 
		Password: password,
	} 

	adminToken := GetAccessToken(credentials)

	return adminToken.AccessToken
}

func UserToKeycloakUser(user User) KeycloakUser {
	keycloakUser := KeycloakUser{
		Username: user.Username, 
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Enabled: true,
	}
	return keycloakUser
}

//Extract a specific claim for the jwt token
func ExtractClaims(tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, nil)
//	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//		return []byte("PUBLIC KEY"), nil
//	})
	if err != nil {
		log.Info("Error while getting claims", "error", err)
	}
	return claims
}

func AddValueToList(values []byte, newKey string, newValue string) []byte {
	jsonBody := map[string]string{}
	err := json.Unmarshal(values, &jsonBody)
	if err != nil {
		log.Info("Error while adding claims", "error", err)
		panic(err)
	}
	jsonBody[newKey] = newValue
	newList, _ := json.Marshal(jsonBody)
	return newList
}
