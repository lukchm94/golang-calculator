package main

import (
	"fmt"
	"net/http"
	"os"

	"app/internal/application"
	authService "app/internal/application/auth"
	calculatorApplication "app/internal/application/calculator"
	userService "app/internal/application/user"
	httpInfra "app/internal/infrastructure/http"
	"app/internal/infrastructure/http/controllers"
)

func main() {
	app := NewApp()

	// Services initialization
	jwtService := authService.NewJwtAuthService(app.Logger, app.JwtConfig)
	healthService := application.NewHealthService(app.Logger)
	calculatorService := calculatorApplication.NewCalculatorService(app.Logger, app.CalcRepo)
	userService := userService.NewUserService(app.Logger, app.UserRepo)
	userController := controllers.NewUserController(app.Logger, userService, jwtService)

	// Controllers initialization
	calculatorController := controllers.NewCalculatorController(app.Logger, calculatorService)

	// HTTP Handlers initialization
	healthHandler := httpInfra.NewHealthHandler(app.Logger, healthService, userService)
	calculatorHandler := httpInfra.NewCalculatorHandler(app.Logger, calculatorController)
	userHandler := httpInfra.NewUserHandler(app.Logger, userController)
	// Router initialization
	adminRouter := httpInfra.NewAdminRouter(app.Logger, jwtService)
	router := httpInfra.NewRouter(healthHandler, calculatorHandler, userHandler, adminRouter, jwtService)

	// Consistent structured logging
	app.Logger.Info("Server starting", "port", app.Config.Port)

	addr := ":" + app.Config.Port // Result: ":8080"

	app.Logger.Info("Starting server",
		"port", app.Config.Port,
		"health_url", fmt.Sprintf("http://localhost:%s/health", app.Config.Port),
	)

	// Start the server
	err := http.ListenAndServe(addr, router)

	addr = ":" + app.Config.Port

	if err != nil {
		app.Logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
