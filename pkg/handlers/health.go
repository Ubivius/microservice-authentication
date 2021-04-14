package handlers

import (
	"net/http"

	"github.com/Ubivius/microservice-authentication/pkg/data"
)

// LivenessCheck determine when the application needs to be restarted
func (authHandler *AuthHandler) LivenessCheck(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
}

//ReadinessCheck verifies that the application is ready to accept requests
func (authHandler *AuthHandler) ReadinessCheck(responseWriter http.ResponseWriter, request *http.Request) {
	readinessProbeKeycloak := data.KeycloakPath + "/auth/realms/master"
	_, err := http.Get(readinessProbeKeycloak)

	if err != nil {
		log.Error(err, "Keycloak unavailable")
		http.Error(responseWriter, "Keycloak unavailable", http.StatusServiceUnavailable)
		return
	}

	readinessProbeMicroserviceUser := data.MicroserviceUserPath + "/health/ready"
	_, err = http.Get(readinessProbeMicroserviceUser)

	if err != nil {
		log.Error(err, "Microservice-user unavailable")
		http.Error(responseWriter, "Microservice-user unavailable", http.StatusServiceUnavailable)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}
