package sqsListener

import (
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsListener struct {
	logger    *slog.Logger
	sqsClient *sqs.Client
	queueUrl  string
}

func NewSqsListener(
	sqsClient *sqs.Client,
	logger *slog.Logger,
	queueUrl string,
) (*SqsListener, error) {

	logger.Info("Initializing SqsListener", "queueUrl", queueUrl)

	return &SqsListener{
		logger:    logger,
		sqsClient: sqsClient,
		queueUrl:  queueUrl,
	}, nil
}

func (l *SqsListener) Listen() {
	l.logger.Info("Starting SQS listener", "queueUrl", l.queueUrl)

	// Implement the logic to listen to the SQS queue and process messages
}
