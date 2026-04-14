package main

import (
	"app/cmd/config"
	dynamoRepo "app/internal/infrastructure/dynamodb/reposiotories"
	eventBridge "app/internal/infrastructure/event_bridge"
	eventBridgeRepo "app/internal/infrastructure/event_bridge/repo"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	dynamodb "app/internal/infrastructure/dynamodb"
	postgres "app/internal/infrastructure/postgres"
	postgresModels "app/internal/infrastructure/postgres/models"
	postgresRepo "app/internal/infrastructure/postgres/repo"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"gorm.io/gorm"
)

type Services struct {
	Logger         *slog.Logger
	CalcRepo       *dynamoRepo.CalculationsDynamoRepository
	Config         config.Configs
	JwtConfig      *config.JwtConfig
	UserRepo       *postgresRepo.UserRepository
	EventPublisher *eventBridgeRepo.EventPublisher
	Context        *context.Context
	Stop           context.CancelFunc
}

func NewApp() *Services {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	appConfig := config.LoadConfigs(logger)

	jwtConfig := buildJwtConfig(appConfig)

	logger.Info("App configuration loaded", "config", appConfig)

	dbClient, err := buildDynamoDbClient(ctx, logger, appConfig)

	if err != nil {
		return nil
	}

	postgresClient, err := buildPostgresClient(logger, appConfig)

	if err != nil {
		logger.Error("Failed to build Postgres client", "error", err)
		return nil
	}

	userRepo := postgresRepo.NewUserRepository(postgresClient, logger)

	calcRepo, err := dynamoRepo.NewCalculationsDynamoRepository(dbClient, logger, appConfig.AwsConfig)

	if err != nil {
		logger.Error("Failed to initialize repository", "error", err)
		os.Exit(1)
	}

	eventBridgeClient, err := buildEventBridgeClient(ctx, logger, appConfig)

	if err != nil {
		logger.Error("Failed to build EventBridge client", "error", err)
		os.Exit(1)
	}

	eventPublisher, err := eventBridgeRepo.NewEventPublisher(eventBridgeClient, logger, appConfig.AwsConfig.EventBus)

	if err != nil {
		logger.Error("Failed to initialize EventBridge repository", "error", err)
		os.Exit(1)
	}
	return &Services{
		Logger:         logger,
		CalcRepo:       calcRepo,
		Config:         appConfig,
		UserRepo:       userRepo,
		JwtConfig:      jwtConfig,
		EventPublisher: eventPublisher,
		Context:        &ctx,
		Stop:           stop,
	}
}

func buildEventBridgeClient(ctx context.Context, logger *slog.Logger, config config.Configs) (*eventbridge.Client, error) {
	eventBridgeCfg, err := eventBridge.LoadEventBridgeConfig(ctx, config.AwsDefaultRegion)

	if err != nil {
		logger.Error("Failed to load EventBridge config", "error", err)
		return nil, err
	}

	eventBridgeClient, err := eventBridge.NewEventBridgeClient(ctx, eventBridge.EventBridgeConfig{
		Config:   eventBridgeCfg,
		Endpoint: config.LocalstackEndpointUrl,
	}, logger)

	if err != nil {
		logger.Error("Failed to create EventBridge client", "error", err)
		return nil, err
	}

	return eventBridgeClient.Client, nil

}

func getDynamoDbCfg(ctx context.Context, logger *slog.Logger, config config.Configs) (aws.Config, error) {
	cfg, err := dynamodb.LoadDynamoDbConfig(ctx, config.AwsDefaultRegion)
	if err != nil {
		logger.Error("Failed to load DynamoDB config", "error", err)
	}

	return cfg, nil
}

func buildDynamoDbClient(ctx context.Context, logger *slog.Logger, config config.Configs) (*dynamodb.DynamoDbClient, error) {
	cfg, err := getDynamoDbCfg(ctx, logger, config)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewDynamoDBClient(ctx, dynamodb.DynamoDbConfig{
		Config:   cfg,
		Endpoint: config.LocalstackEndpointUrl,
	}, logger)
}

func buildPostgresClient(logger *slog.Logger, config config.Configs) (*gorm.DB, error) {
	postgresConfig := postgres.PostgresConfig{
		Host:     config.PostgresHost,
		Port:     config.PostgresPort,
		User:     config.PostgresUser,
		Password: config.PostgresPassword,
		DbName:   config.PostgresDb,
	}
	db, err := postgres.NewGormClient(postgresConfig, logger)

	if err != nil {
		logger.Error("Failed to build Postgres client", "error", err)

		return nil, postgres.ErrPostgresInit

	}

	logger.Info("Initialised Postgres client", "Host", postgresConfig.Host, "Port", postgresConfig.Port, "DbName", postgresConfig.DbName)

	if err := postgres.InitPostgresTables(db, logger); err != nil {
		logger.Error("Failed to initialize Postgres tables", "error", err)
		return nil, err
	}

	logger.Info("Initialised Postgres tables", "Tables", []string{postgresModels.UserPostgres{}.TableName()})

	return db, nil
}

func buildJwtConfig(c config.Configs) *config.JwtConfig {
	issuer := config.Issuer(c.AppIssuer)
	if !issuer.IsValid() {
		issuer = config.AppIssuer
	}
	expTime := config.FromStringToTimeDuration(c.JwtExpirationTime)

	return &config.JwtConfig{
		SecretKey:      c.JwtSecretKey,
		Issuer:         issuer,
		ExpirationTime: expTime,
	}
}
