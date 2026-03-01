package controllers

import (
	userService "app/internal/application/user"
	userDomain "app/internal/domain/user"
	reqErr "app/internal/infrastructure/http/errors"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
)

type UserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
}

type LoginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserController struct {
	logger  *slog.Logger
	service *userService.UserService
}

func NewUserController(logger *slog.Logger, service *userService.UserService) *UserController {
	return &UserController{logger: logger, service: service}
}

func (c *UserController) Register(ctx context.Context, r *http.Request) (*userDomain.User, error) {
	validReq, err := c.validateRegisterReq(r)

	if err != nil {
		c.logger.Error("Request validation failed", "error", err)

		return nil, err
	}

	result, err := c.service.Register(validReq)

	if err != nil {
		c.logger.Error("Failed to register user", "error", err)

		return nil, err
	}

	if result != nil {
		c.logger.Info("Successfully registered user", "email", result.Email, "firsName", result.FirstName, "lastName", result.LastName)
	}

	return result, nil
}

func (c *UserController) Login(ctx context.Context, r *http.Request) (*userDomain.User, error) {
	c.logger.Info("Handling Login")

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, reqErr.InvalidRequestError{}
	}

	if req.Email == nil {
		return nil, reqErr.MissingFieldError{FieldName: "email"}
	}

	err := c.isValidEmail(*req.Email)

	if err != nil {
		return nil, err
	}
	loginInput := userService.LoginInput{
		Username: strings.TrimSpace(*req.Email),
		Password: strings.TrimSpace(*req.Password),
	}

	user, err := c.service.Login(loginInput)

	if err != nil {
		c.logger.Error("Failed to login user", "error", err)

		if errors.Is(err, userDomain.ErrUserNotFound) {
			var errDetail = fmt.Sprintf("username: %s", loginInput.Username)

			return nil, reqErr.UserNotFoundError{Details: errDetail}
		}

		if errors.Is(err, userDomain.ErrInvalidCredentials) {
			var errDetail = fmt.Sprintf("invalid credentials for username: %s", loginInput.Username)

			return nil, reqErr.InvalidCredentialsError{Details: errDetail}
		}

		return nil, err

	}

	c.logger.Info("Successfully found user", "email", user.Email, "firstName", user.FirstName, "lastName", user.LastName)

	return user, nil
}

func (c *UserController) validateRegisterReq(r *http.Request) (userService.RegisterInput, error) {
	c.logger.Info("Handling register request")

	if r.Method != http.MethodPost {
		c.logger.Error("Invalid HTTP method. Expected: POST received", "method", r.Method)

		return userService.RegisterInput{}, reqErr.InvalidRequestMethodError{Method: http.MethodPost}
	}

	req, err := c.validateRegisterPayload(r)

	if err != nil {
		c.logger.Error("Invalid request", "error", err)

		return userService.RegisterInput{}, err
	}

	return *req, err
}

func (c *UserController) validateRegisterPayload(r *http.Request) (*userService.RegisterInput, error) {
	var req UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, reqErr.InvalidRequestError{}
	}

	if req.Email == nil {
		return nil, reqErr.MissingFieldError{FieldName: "email"}
	}
	err := c.isValidEmail(*req.Email)

	if err != nil {
		return nil, err
	}

	if req.FirstName == nil {
		return nil, reqErr.MissingFieldError{FieldName: "firstName"}
	}

	if req.LastName == nil {
		return nil, reqErr.MissingFieldError{FieldName: "lastName"}
	}

	if req.Password == nil {
		return nil, reqErr.MissingFieldError{FieldName: "password"}
	}

	return &userService.RegisterInput{
		FirstName: *req.FirstName,
		LastName:  *req.LastName,
		Email:     *req.Email,
		Password:  *req.Password,
	}, nil
}

func (c *UserController) isValidEmail(email string) error {
	address, err := mail.ParseAddress(email)
	if err != nil {
		var errDetail = fmt.Sprintf("Invalid email: %s", email)
		return reqErr.InvalidRequestError{Details: errDetail}
	}
	c.logger.Info("Validated email address", "email", address.Address)

	return nil
}
