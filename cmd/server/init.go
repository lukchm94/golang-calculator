package main

import (
	"app/cmd/config"
	dynamodbModels "app/internal/infrastructure/dynamodb/models"
	dynamoRepo "app/internal/infrastructure/dynamodb/reposiotories"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	dynamodb "app/internal/infrastructure/dynamodb"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type Services struct {
	Logger   *slog.Logger
	CalcRepo *dynamoRepo.CalculationsDynamoRepository
	Config   config.Configs
}

func NewApp() *Services {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	appConfig := config.LoadConfigs()

	logger.Info("App configuration loaded", "config", appConfig)

	dbClient, err := buildDynamoDbClient(ctx, logger, appConfig)

	if err != nil {
		return nil
	}

	initDb := dynamodb.NewInitDynamoDb(logger, dbClient)

	if err := initDb.EnsureTablesExist(ctx); err != nil {
		logger.Error("Failed to ensure DynamoDB tables exist", "error", err)
		return nil
	}

	if err != nil {
		logger.Error("Failed to build DynamoDB client", "error", err)
		return nil
	}

	calcRepo, err := dynamoRepo.NewCalculationsDynamoRepository(dbClient, logger, "Calculations")

	if err != nil {
		logger.Error("Failed to initialize repository", "error", err)
		os.Exit(1)
	}

	return &Services{
		Logger:   logger,
		CalcRepo: calcRepo,
		Config:   appConfig,
	}
}

func getAwsConfig(ctx context.Context, logger *slog.Logger, config config.Configs) (aws.Config, error) {
	cfg, err := dynamodb.LoadAWSConfig(ctx, config.AwsDefaultRegion)
	if err != nil {
		logger.Error("Failed to load AWS config", "error", err)
	}

	return cfg, nil
}

func buildDynamoDbClient(ctx context.Context, logger *slog.Logger, config config.Configs) (*dynamodb.DynamoDbClient, error) {
	cfg, err := getAwsConfig(ctx, logger, config)
	if err != nil {
		return nil, err
	}

	tables := dynamoTablesToRegister()
	logger.Info("DynamoDB tables to register", "tables", tables)

	return dynamodb.NewDynamoDBClient(ctx, dynamodb.DynamoDbConfig{
		Config:   cfg,
		Endpoint: config.LocalstackEndpointUrl,
		Tables:   dynamoTablesToRegister(),
	}, logger)
}

func dynamoTablesToRegister() []string {
	return []string{dynamodbModels.TABLE_CALCULATIONS}
}
