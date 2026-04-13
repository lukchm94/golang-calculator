package dynamodb

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDbConfig struct {
	Config   aws.Config
	Endpoint string
}

type DynamoDbClient struct {
	logger *slog.Logger
	Client *dynamodb.Client
}

func NewDynamoDBClient(context context.Context, input DynamoDbConfig, logger *slog.Logger) (*DynamoDbClient, error) {
	logger.Info("Creating DynamoDB client", "endpoint", input.Endpoint)

	sdkClient := dynamodb.NewFromConfig(input.Config, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(input.Endpoint)
	})

	return &DynamoDbClient{
		logger: logger,
		Client: sdkClient,
	}, nil
}

func LoadAWSConfig(ctx context.Context, region string) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
}
