package main

import (
	"fmt"
	"net/http"
	"os"

	"app/cmd/config"
	"app/internal/application"
	authService "app/internal/application/auth"
	calculatorApplication "app/internal/application/calculator"
	eventsDispatcher "app/internal/application/events"
	userService "app/internal/application/user"
	userLoginHandler "app/internal/application/user/handlers"
	httpInfra "app/internal/infrastructure/http"
	"app/internal/infrastructure/http/controllers"
	"app/internal/infrastructure/sqs"
	sqsListener "app/internal/infrastructure/sqs/repo"
)

func main() {
	app := NewApp()
	defer app.Stop()

	// Services initialization
	jwtService := authService.NewJwtAuthService(app.Logger, app.JwtConfig)
	healthService := application.NewHealthService(app.Logger)
	calculatorService := calculatorApplication.NewCalculatorService(app.Logger, app.CalcRepo)
	userMapper := userService.NewUserMapper(app.Logger)
	userService := userService.NewUserService(app.Logger, app.UserRepo, app.EventPublisher, userMapper)
	userController := controllers.NewUserController(app.Logger, userService, jwtService)

	// Events dispatcher initialization
	sqsRegion, err := sqs.LoadSqsConfig(*app.Context, string(app.Config.AwsDefaultRegion))
	if err != nil {
		app.Logger.Error("Failed to load SQS config", "error", err)
		os.Exit(1)
	}
	sqsConfig := sqs.SqsConfig{
		Config:   sqsRegion,
		Endpoint: app.Config.LocalstackEndpointUrl,
	}
	sqsClient, err := sqs.NewSqsClient(*app.Context, sqsConfig, app.Logger)
	if err != nil {
		app.Logger.Error("Failed to create SQS client", "error", err)
		os.Exit(1)
	}
	userLoginHandler := userLoginHandler.NewUserLoginHandler(app.Logger)
	sqsDispatcher := eventsDispatcher.NewSqsDispatcher(app.Logger, userLoginHandler, userMapper)

	sqsListener, err := sqsListener.NewSqsListener(sqsClient.Client, app.Logger, sqsDispatcher, string(config.SqsQueueUrlPrefixLocal)+app.Config.AwsConfig.MainQueueName)
	if err != nil {
		app.Logger.Error("Failed to initialize SQS listener", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := sqsListener.Listen(*app.Context); err != nil {
			app.Logger.Error("SQS listener stopped with error", "error", err)
			app.Stop()
		}
	}()

	// Controllers initialization
	calculatorController := controllers.NewCalculatorController(app.Logger, calculatorService)

	// HTTP Handlers initialization
	healthHandler := httpInfra.NewHealthHandler(app.Logger, healthService, userService)
	calculatorHandler := httpInfra.NewCalculatorHandler(app.Logger, calculatorController)
	userHandler := httpInfra.NewUserHandler(app.Logger, userController)

	// Router initialization
	adminRouter := httpInfra.NewAdminRouter(app.Logger, jwtService)
	router := httpInfra.NewRouter(
		healthHandler,
		calculatorHandler,
		userHandler,
		adminRouter,
		jwtService,
		app.Logger,
	)

	// Consistent structured logging
	app.Logger.Info("Server starting", "port", app.Config.Port)

	addr := ":" + app.Config.Port // Result: ":8080"

	app.Logger.Info("Starting server",
		"port", app.Config.Port,
		"health_url", fmt.Sprintf("http://localhost:%s/health", app.Config.Port),
	)
	app.Logger.Info("Starting SQS listener", "queueUrl", string(config.SqsQueueUrlPrefixLocal)+app.Config.AwsConfig.MainQueueName)

	// Start the server
	err = http.ListenAndServe(addr, router)

	if err != nil {
		app.Logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
