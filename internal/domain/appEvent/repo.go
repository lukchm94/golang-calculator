package appEvent

import "context"

type EventPublisher interface {
	PublishEvent(ctx context.Context, event PublishingEvent) error
}
