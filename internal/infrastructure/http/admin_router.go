package httpInfra

import (
	"log/slog"
	"net/http"

	authService "app/internal/application/auth"
	"app/internal/infrastructure/http/middleware"
)

type AdminRouter struct {
	mux        *http.ServeMux
	logger     *slog.Logger
	jwtService *authService.JwtAuthService
}

func NewAdminRouter(logger *slog.Logger, jwtService *authService.JwtAuthService) *AdminRouter {
	return &AdminRouter{
		mux:        http.NewServeMux(),
		logger:     logger,
		jwtService: jwtService,
	}
}

func (ar *AdminRouter) HandleFunc(path string, handler http.HandlerFunc) {
	chain := middleware.AuthMiddleware(ar.logger, ar.jwtService)(
		middleware.AdminOnlyMiddleware(ar.logger)(handler),
	)
	ar.mux.Handle(path, chain)
}

func (ar *AdminRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ar.logger.Info("AdminRouter received request", "path", r.RequestURI, "method", r.Method)
	ar.mux.ServeHTTP(w, r)
}
