package appEvent

import "context"

type EventPublisher interface {
	PublishEvent(ctx context.Context, event PublishingEvent) error
}

type SqsListener interface {
	Listen(ctx context.Context) error
}
