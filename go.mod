module github.com/Ubivius/microservice-authentication

go 1.15

require (
	github.com/Ubivius/pkg-telemetry v1.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/stretchr/objx v0.2.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.26.0
	go.opentelemetry.io/otel v1.1.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	sigs.k8s.io/controller-runtime v0.10.2
)
