package httpInfra

import (
	"net/http"
)

func NewRouter(healthHandler *HealthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(string(HealthRoute), healthHandler)

	return mux
}
