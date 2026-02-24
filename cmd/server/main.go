package main

import (
	"fmt"
	"net/http"
	"os"

	"app/internal/application"
	calculatorApplication "app/internal/application/calculator"
	userService "app/internal/application/user"
	httpInfra "app/internal/infrastructure/http"
	"app/internal/infrastructure/http/controllers"
)

func main() {
	app := NewApp()

	// Services initialization
	healthService := application.NewHealthService(app.Logger)
	calculatorService := calculatorApplication.NewCalculatorService(app.Logger, app.CalcRepo)
	userService := userService.NewUserService(app.Logger, app.UserRepo)
	userController := controllers.NewUserController(app.Logger, userService)

	// Controllers initialization
	calculatorController := controllers.NewCalculatorController(app.Logger, calculatorService)

	// HTTP Handlers initialization
	healthHandler := httpInfra.NewHealthHandler(app.Logger, healthService)
	calculatorHandler := httpInfra.NewCalculatorHandler(app.Logger, calculatorController)
	userHandler := httpInfra.NewUserHandler(app.Logger, userController)
	// Router initialization
	router := httpInfra.NewRouter(healthHandler, calculatorHandler, userHandler)

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

	// ListenAndServe blocks until there is an error
	err = http.ListenAndServe(addr, router)

	if err != nil {
		app.Logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
