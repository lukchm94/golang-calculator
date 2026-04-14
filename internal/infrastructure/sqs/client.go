package sqs

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsConfig struct {
	Config   aws.Config
	Endpoint string
}

type SqsClient struct {
	logger *slog.Logger
	Client *sqs.Client
}

func NewSqsClient(context context.Context, input SqsConfig, logger *slog.Logger) (*SqsClient, error) {
	logger.Info("Creating SQS client", "endpoint", input.Endpoint)

	sdkClient := sqs.NewFromConfig(input.Config, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(input.Endpoint)
	})

	return &SqsClient{
		logger: logger,
		Client: sdkClient,
	}, nil
}

func LoadSqsConfig(ctx context.Context, region string) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
}
