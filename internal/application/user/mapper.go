package userService

import (
	userEvents "app/internal/application/user/events"
	"app/internal/domain/appEvent"
	userDomain "app/internal/domain/user"
	"encoding/json"
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

func (m *UserMapper) FromPublishingEventToLoginEvent(event appEvent.PublishingEvent) (*userEvents.LoginEvent, error) {
	m.logger.Debug("Mapping PublishingEvent to LoginEvent", "publishingEvent", event)

	var loginEvent userEvents.LoginEvent
	err := json.Unmarshal([]byte(event.Detail), &loginEvent)

	if err != nil {
		m.logger.Error("Failed to unmarshal PublishingEvent detail into LoginEvent", "error", err, "detail", event.Detail)
		return nil, err
	}

	return &loginEvent, nil
}
