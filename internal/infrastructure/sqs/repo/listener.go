package sqsListener

import (
	eventsDispatcher "app/internal/application/events"
	"app/internal/domain/appEvent"
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SqsListener struct {
	logger           *slog.Logger
	sqsClient        *sqs.Client
	eventsDispatcher *eventsDispatcher.SqsDispatcher
	queueUrl         string
}

func NewSqsListener(
	sqsClient *sqs.Client,
	logger *slog.Logger,
	eventsDispatcher *eventsDispatcher.SqsDispatcher,
	queueUrl string,
) (*SqsListener, error) {

	logger.Info("Initializing SqsListener", "queueUrl", queueUrl)

	return &SqsListener{
		logger:           logger,
		sqsClient:        sqsClient,
		eventsDispatcher: eventsDispatcher,
		queueUrl:         queueUrl,
	}, nil
}

func (l *SqsListener) Listen(ctx context.Context) error {
	l.logger.Info("Starting SQS listener", "queueUrl", l.queueUrl)

	for {
		if err := ctx.Err(); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				l.logger.Info("Stopping SQS listener", "queueUrl", l.queueUrl, "reason", err)
				return nil
			}

			return err
		}

		resp, err := l.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &l.queueUrl,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,
		})
		if err != nil {
			l.logger.Error("Failed to receive messages from SQS", "error", err, "queueUrl", l.queueUrl)
			return err
		}

		if len(resp.Messages) == 0 {
			continue
		}

		for _, message := range resp.Messages {
			if err := l.handleMessage(ctx, message); err != nil {
				l.logger.Error("Failed to process SQS message", "error", err, "messageId", safeString(message.MessageId))
			}
		}
	}
}

func (l *SqsListener) handleMessage(ctx context.Context, message types.Message) error {
	l.logger.Info("Received SQS message", "messageId", safeString(message.MessageId))

	event, err := l.decodePublishingEvent(message)
	if err != nil {
		return err
	}

	if err := l.eventsDispatcher.Dispatch(event); err != nil {
		return err
	}

	return l.deleteMessage(ctx, message)
}

func (l *SqsListener) deleteMessage(ctx context.Context, message types.Message) error {
	if message.ReceiptHandle == nil {
		return errors.New("cannot delete SQS message without receipt handle")
	}

	_, err := l.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &l.queueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		l.logger.Error("Failed to delete SQS message", "error", err, "messageId", safeString(message.MessageId))
		return err
	}

	l.logger.Info("Deleted SQS message after successful handling", "messageId", safeString(message.MessageId))

	return nil
}

type eventBridgeEnvelope struct {
	Source     appEvent.EventSource        `json:"source"`
	DetailType appEvent.AppEventDetailType `json:"detail-type"`
	Detail     json.RawMessage             `json:"detail"`
}

func (l *SqsListener) decodePublishingEvent(message types.Message) (appEvent.PublishingEvent, error) {
	if message.Body == nil {
		return appEvent.PublishingEvent{}, errors.New("received SQS message with empty body")
	}

	var envelope eventBridgeEnvelope
	if err := json.Unmarshal([]byte(*message.Body), &envelope); err != nil {
		l.logger.Error("Failed to decode EventBridge envelope from SQS message", "error", err, "body", *message.Body)
		return appEvent.PublishingEvent{}, err
	}

	return appEvent.PublishingEvent{
		Source:     envelope.Source,
		DetailType: envelope.DetailType,
		Detail:     string(envelope.Detail),
	}, nil
}

func safeString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
