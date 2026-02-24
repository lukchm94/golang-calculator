package controllers

import (
	userService "app/internal/application/user"
	userDomain "app/internal/domain/users"
	"context"
	"errors"
	"log/slog"
	"net/http"
)

type UserController struct {
	logger  *slog.Logger
	service *userService.UserService
}

func NewUserController(logger *slog.Logger, service *userService.UserService) *UserController {
	return &UserController{logger: logger, service: service}
}

func (c *UserController) Register(ctx context.Context, r *http.Request) (userDomain.User, error) {
	return userDomain.User{}, errors.New("not implemented")
}
