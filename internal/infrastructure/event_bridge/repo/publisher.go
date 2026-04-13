package eventBridgeRepo

import (
	"app/internal/domain/appEvent"
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

type EventPublisher struct {
	logger            *slog.Logger
	eventBridgeClient *eventbridge.Client
	eventBusArn       string
}

func NewEventPublisher(
	eventBridgeClient *eventbridge.Client,
	logger *slog.Logger,
	eventBusArn string,
) (*EventPublisher, error) {

	logger.Info("Initializing EventBridgeRepository", "eventBusArn", eventBusArn)

	return &EventPublisher{
		logger:            logger,
		eventBridgeClient: eventBridgeClient,
		eventBusArn:       eventBusArn,
	}, nil
}

func (r *EventPublisher) PublishEvent(ctx context.Context, event appEvent.PublishingEvent) error {

	r.logger.Debug("Publishing event to EventBridge", "event", event)

	putEventsIpnut := &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Source:       aws.String(event.Source.String()),
				DetailType:   aws.String(event.DetailType.String()),
				Detail:       aws.String(event.Detail),
				EventBusName: aws.String(r.eventBusArn),
			},
		},
	}
	resp, err := r.eventBridgeClient.PutEvents(ctx, putEventsIpnut)

	r.logger.Debug("EventBridge PutEvents response", "response", resp)

	if err != nil {
		r.logger.Error("Failed to publish event to EventBridge", "error", err)
		return err
	}

	return nil
}
