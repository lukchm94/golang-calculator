package eventsDispatcher

import (
	userService "app/internal/application/user"
	userLoginHandler "app/internal/application/user/handlers"
	"app/internal/domain/appEvent"
	"log/slog"
)

type SqsDispatcher struct {
	logger       *slog.Logger
	loginHandler *userLoginHandler.UserLoginHandler
	userMapper   *userService.UserMapper
}

func NewSqsDispatcher(
	logger *slog.Logger,
	loginHandler *userLoginHandler.UserLoginHandler,
	userMapper *userService.UserMapper,
) *SqsDispatcher {
	return &SqsDispatcher{
		logger:       logger,
		loginHandler: loginHandler,
		userMapper:   userMapper,
	}
}

func (d *SqsDispatcher) Dispatch(event appEvent.PublishingEvent) error {
	d.logger.Info("Dispatching message from SQS", "message", event.Detail)

	switch event.DetailType {
	case appEvent.Login:
		return d.handleLoginEvent(event)

	default:
		d.logger.Warn("Received unknown event type", "eventType", event.DetailType)

		return NotImplementedError
	}

}

func (d *SqsDispatcher) handleLoginEvent(event appEvent.PublishingEvent) error {
	d.logger.Info("Handling user login event", "eventDetail", event.Detail)

	loginEvent, err := d.userMapper.FromPublishingEventToLoginEvent(event)

	if err != nil {
		d.logger.Error("Failed to map PublishingEvent to LoginEvent", "error", err)
		return err
	}

	return d.loginHandler.HandleLogin(*loginEvent)
}
