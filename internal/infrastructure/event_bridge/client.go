package eventBridge

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

type EventBridgeConfig struct {
	Config   aws.Config
	Endpoint string
}

type EventBridgeClient struct {
	logger *slog.Logger
	Client *eventbridge.Client
}

func NewEventBridgeClient(context context.Context, input EventBridgeConfig, logger *slog.Logger) (*EventBridgeClient, error) {
	logger.Info("Creating EventBridge client", "endpoint", input.Endpoint)

	sdkClient := eventbridge.NewFromConfig(input.Config, func(o *eventbridge.Options) {
		o.BaseEndpoint = aws.String(input.Endpoint)
	})

	return &EventBridgeClient{
		logger: logger,
		Client: sdkClient,
	}, nil
}

func LoadEventBridgeConfig(ctx context.Context, region string) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
}
