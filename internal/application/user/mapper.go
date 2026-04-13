package userService

import (
	userEvents "app/internal/application/user/events"
	"app/internal/domain/appEvent"
	userDomain "app/internal/domain/user"
	"log/slog"
	"time"
)

type UserMapper struct {
	logger *slog.Logger
}

func NewUserMapper(logger *slog.Logger) *UserMapper {
	return &UserMapper{
		logger: logger,
	}
}

func (m *UserMapper) FromUserDomainToLoginEvent(user *userDomain.User, timestamp time.Time) *userEvents.LoginEvent {
	m.logger.Debug("Mapping User domain model to LoginEvent", "userID", user.ID)

	return &userEvents.LoginEvent{
		UserID:    user.ID,
		Timestamp: timestamp,
		Role:      user.Role,
	}
}

func (m *UserMapper) FromLoginEventToPublishingEvent(event *userEvents.LoginEvent) appEvent.PublishingEvent {
	publishingEvent := appEvent.PublishingEvent{
		Source:     appEvent.CalculatorApp,
		DetailType: appEvent.Login,
		Detail:     event.JSON(),
	}

	m.logger.Debug("Mapped LoginEvent to PublishingEvent", "publishingEvent", publishingEvent)

	return publishingEvent
}
