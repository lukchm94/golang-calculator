package httpInfra

import (
	authService "app/internal/application/auth"
	"log/slog"
	"net/http"
	"os"
)

func NewRouter(
	healthHandler *HealthHandler,
	calculatorHandler *CalculatorHandler,
	userHandler *UserHandler,
	adminRouter *AdminRouter,
	jwtService *authService.JwtAuthService,
) http.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	mux := http.NewServeMux()

	// Public routes
	healthHandler.RegisterRoutes(mux)
	calculatorHandler.RegisterRoutes(mux)
	userHandler.RegisterRoutes(mux)

	// Admin routes
	healthHandler.RegisterAdminRoutes(adminRouter)

	logger.Info("Registering admin router at path", "path", string(AdminRoute))
	mux.Handle(string(AdminRoute), adminRouter)

	return mux
}
