package userLoginHandler

import (
	userEvents "app/internal/application/user/events"
	"log/slog"
)

type UserLoginHandler struct {
	logger *slog.Logger
}

func NewUserLoginHandler(logger *slog.Logger) *UserLoginHandler {
	return &UserLoginHandler{
		logger: logger,
	}
}

func (h *UserLoginHandler) HandleLogin(loginEvent userEvents.LoginEvent) error {
	h.logger.Info("Handling user login", "userID", loginEvent.UserID)
	// Implement the logic to handle user login, such as generating tokens, updating last login time, etc.
	return nil
}
