package router

import (
	"log"
	"net/http"

	"github.com/Ubivius/microservice-authentication/pkg/handlers"
	"github.com/gorilla/mux"
)

// Mux route handling with gorilla/mux
func New(authHandler *handlers.AuthHandler, logger *log.Logger) *mux.Router {
	router := mux.NewRouter()

	// Post router
	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/signin", authHandler.SignIn)
	postRouter.HandleFunc("/signup", authHandler.SignUp)
	//postRouter.Use(authHandler.MiddlewareAuthValidation)

	return router
}
