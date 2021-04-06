package handlers

import (
	"net/http"

	"github.com/Ubivius/microservice-authentication/pkg/data"
)

// LivenessCheck determine when the application needs to be restarted
func (authHandler *AuthHandler) LivenessCheck(responseWriter http.ResponseWriter, request *http.Request) {
	log.Info("LivenessCheck")
	responseWriter.WriteHeader(http.StatusOK)
}

//ReadinessCheck verifies that the application is ready to accept requests
func (authHandler *AuthHandler) ReadinessCheck(responseWriter http.ResponseWriter, request *http.Request) {
	log.Info("ReadinessCheck")

	readinessProbeKeycloak := data.KeycloakPath + "/auth/realms/master"
	readinessProbeMicroserviceUser := data.MicroserviceUserPath + "/health/ready"

	_, errKeycloak := http.Get(readinessProbeKeycloak)
	_, errMicroserviceUser := http.Get(readinessProbeMicroserviceUser)

	if errKeycloak != nil {
		log.Error(errKeycloak, "Keycloak unavailable")
		http.Error(responseWriter, "Keycloak unavailable", http.StatusServiceUnavailable)
		return
	}

	if errMicroserviceUser != nil {
		log.Error(errMicroserviceUser, "Microservice-user unavailable")
		http.Error(responseWriter, "Microservice-user unavailable", http.StatusServiceUnavailable)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}
