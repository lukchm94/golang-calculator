package sqs

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSConfig struct {
	Config   aws.Config
	Endpoint string
}

type SQSClient struct {
	logger *slog.Logger
	Client *sqs.Client
}

func NewSQSClient(context context.Context, input SQSConfig, logger *slog.Logger) (*SQSClient, error) {
	logger.Info("Creating SQS client", "endpoint", input.Endpoint)

	sdkClient := sqs.NewFromConfig(input.Config, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(input.Endpoint)
	})

	return &SQSClient{
		logger: logger,
		Client: sdkClient,
	}, nil
}

func LoadSQSConfig(ctx context.Context, region string) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
}
