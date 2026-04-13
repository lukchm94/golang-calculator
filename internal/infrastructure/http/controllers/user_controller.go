package controllers

import (
	authService "app/internal/application/auth"
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
	Role      *string `json:"role"`
}

type LoginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserController struct {
	logger     *slog.Logger
	service    *userService.UserService
	jwtService *authService.JwtAuthService
}

func NewUserController(logger *slog.Logger, service *userService.UserService, jwtService *authService.JwtAuthService) *UserController {
	return &UserController{logger: logger, service: service, jwtService: jwtService}
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

		if errors.Is(err, userDomain.ErrUserAlreadyExists) {
			return nil, reqErr.UserAlreadyExistsError{Details: validReq.Email}
		}

		return nil, err
	}

	if result != nil {
		c.logger.Info("Successfully registered user", "email", result.Email, "firsName", result.FirstName, "lastName", result.LastName)
	}

	return result, nil
}

func (c *UserController) Login(ctx context.Context, r *http.Request) (*UserLoginResponse, error) {
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

	if req.Password == nil || strings.TrimSpace(*req.Password) == "" {
		return nil, reqErr.MissingFieldError{FieldName: "password"}
	}

	loginInput := userService.LoginInput{
		Username: strings.TrimSpace(*req.Email),
		Password: strings.TrimSpace(*req.Password),
	}

	user, err := c.service.Login(ctx, loginInput)

	if err != nil {
		return c.handleLoginError(loginInput, err)
	}

	c.logger.Info("Successfully found user", "email", user.Email, "firstName", user.FirstName, "lastName", user.LastName)

	return c.generateUserWithToken(user)
}

func (c *UserController) generateUserWithToken(user *userDomain.User) (*UserLoginResponse, error) {
	resp := &UserLoginResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	if user.Role == userDomain.Admin {
		token, err := c.jwtService.GenerateToken(authService.JwtLoginInput{
			UserID: user.ID,
			Role:   user.Role,
		})
		if err != nil {
			return nil, err
		}
		resp.Token = c.maskToken(token)
		c.logger.Info("Generated JWT token for user", "email", user.Email, "maskedToken", resp.Token)
	}

	return resp, nil
}
func (c *UserController) maskToken(token string) string {
	if len(token) <= 10 {
		c.logger.Warn("Token length is too short to mask properly", "tokenLength", len(token))
		return "****"
	}

	masked := token[:5] + "****" + token[len(token)-5:]
	c.logger.Debug("Masked token for logging", "maskedToken", masked)

	return masked
}

func (c *UserController) handleLoginError(loginInput userService.LoginInput, err error) (*UserLoginResponse, error) {
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

	if req.Role == nil {
		return nil, reqErr.MissingFieldError{FieldName: "role"}
	}

	role := userDomain.Role(*req.Role)

	if !role.IsValid() {
		var errDetail = fmt.Sprintf("Invalid role: %s", *req.Role)
		return nil, reqErr.InvalidRequestError{Details: errDetail}
	}

	return &userService.RegisterInput{
		FirstName: *req.FirstName,
		LastName:  *req.LastName,
		Email:     *req.Email,
		Password:  *req.Password,
		Role:      *req.Role,
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
