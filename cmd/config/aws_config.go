package config

import "log/slog"

type AwsPrefix string

const (
	DevDynamoDbPrefix   AwsPrefix = "dev-"
	StageDynamoDbPrefix AwsPrefix = "stage-"
	ProdDynamoDbPrefix  AwsPrefix = "prod-"
)

type EventBridge string

const (
	CalculatorEvents EventBridge = "calculator-events"
)

type EventBridgeArnPrefix string

const (
	EventBridgeArnPrefixLocal EventBridgeArnPrefix = "arn:aws:events:eu-central-1:000000000000:event-bus/"
)

type SqsQueue string

const (
	MainQueueName   SqsQueue = "calculator-main-queue"
	MainDlqName     SqsQueue = "calculator-dlq"
	DeliveryDlqName SqsQueue = "calculator-delivery-dlq"
)

type SqsQueueUrlPrefix string

const (
	SqsQueueUrlPrefixLocal SqsQueueUrlPrefix = "http://sqs.eu-central-1.localhost.localstack.cloud:4566/000000000000/"
)

type AwsConfig struct {
	Prefix          AwsPrefix
	EventBus        string
	MainQueueName   string
	MainDlqName     string
	DeliveryDlqName string
}

func GetAwsConfig(env ValidEnvironments, logger *slog.Logger) AwsConfig {
	prefix := resolveDynamoDbPrefix(env)

	aws := AwsConfig{
		Prefix:          prefix,
		EventBus:        buildEventBridgeName(prefix, CalculatorEvents),
		MainQueueName:   buildSqsQueueName(prefix, MainQueueName),
		MainDlqName:     buildSqsQueueName(prefix, MainDlqName),
		DeliveryDlqName: buildSqsQueueName(prefix, DeliveryDlqName),
	}
	logger.Info("AWS configuration built", "awsConfig", aws)

	return aws
}

func resolveDynamoDbPrefix(env ValidEnvironments) AwsPrefix {
	switch env {
	case DevEnvironment:
		return DevDynamoDbPrefix
	case StageEnvironment:
		return StageDynamoDbPrefix
	case ProdEnvironment:
		return ProdDynamoDbPrefix
	default:
		return DevDynamoDbPrefix
	}
}

func buildEventBridgeName(prefix AwsPrefix, eventBridge EventBridge) string {
	return string(prefix) + string(eventBridge)
}

func buildSqsQueueName(prefix AwsPrefix, queueName SqsQueue) string {
	return string(prefix) + string(queueName)
}
