# microservice-authentication
Microservice to Sign in and Sign up users to [Keycloak](https://github.com/Ubivius/deploy-keycloak)

## Authentication endpoints

`GET` `/health/live` Returns a Status OK when live.

`GET` `/health/ready` Returns a Status OK when ready or an error when dependencies are not available.

`POST` `/signin` Signs in a user with their credentials and returns an access token along with the user's information.

###### Data Params

```json
{
  "username": "string, required",
  "password": "string, required",
}
```

`POST` `/signup` Add a user with their credentials in [Keycloak](https://github.com/Ubivius/deploy-keycloak) and create a new user with [microservice-user](https://github.com/Ubivius/microservice-user). 

###### Data Params

```json
{
  "username":    "string, required",
  "password":    "string, required",
  "firstname":   "string, required",
  "lastname":    "string, required",
  "email":       "string, required",
  "dateofbirth": "string, required",
}
```
