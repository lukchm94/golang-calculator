package userService

import (
	userDomain "app/internal/domain/user"
	"app/internal/utils"
	"log/slog"
	"strings"
)

type UserService struct {
	logger *slog.Logger
	repo   userDomain.UserRepository
}

func NewUserService(logger *slog.Logger, repo userDomain.UserRepository) *UserService {
	return &UserService{
		logger: logger,
		repo:   repo,
	}
}

func (s *UserService) Register(input RegisterInput) (*userDomain.User, error) {
	s.logger.Info("Registering new user", "email", input.Email)
	s.logger.Info("Validating role for user registration", "role", input.Role)

	role := userDomain.Role(input.Role)

	s.logger.Info("Parsed role for user registration", "role", role)

	if !role.IsValid() {
		s.logger.Info("Invalid role provided for user registration", "role", input.Role)
		return nil, userDomain.ErrInvalidRole
	}

	user := &userDomain.User{
		ID:        s.repo.GenerateUserID(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Role:      role,
	}

	err := user.SetPassword(input.Password)
	if err != nil {
		s.logger.Error("Failed to encode password", "error", err)
		return nil, err
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		s.logger.Error("Failed to create user", "error", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(input LoginInput) (*userDomain.User, error) {
	user, err := s.findWithUsername(input.Username)

	userVerificationError := s.verifyUser(user, err)

	if userVerificationError != nil {
		return nil, userVerificationError
	}

	err = user.VerifyPassword(input.Password)

	if err != nil {
		s.logger.Info("Invalid password for user", "username", input.Username)
		return nil, userDomain.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) findWithUsername(username string) (*userDomain.User, error) {

	if strings.Contains(username, utils.EmailCharacter) {
		return s.findWithEmail(username)
	}

	return s.findWithID(username)
}

func (s *UserService) findWithEmail(email string) (*userDomain.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		s.logger.Error("Failed to get user by email", "error", err)
		return nil, err
	}

	if user == nil {
		s.logger.Info("User not found with email", "email", email)
		return nil, nil
	}

	return user, nil
}

func (s *UserService) findWithID(id string) (*userDomain.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		s.logger.Error("Failed to get user by ID", "error", err)
		return nil, err
	}

	if user == nil {
		s.logger.Info("User not found with ID", "id", id)
		return nil, nil
	}

	return user, nil
}

func (s *UserService) verifyUser(user *userDomain.User, err error) error {
	if err != nil {
		s.logger.Error("Failed to get user", "error", err)
		return err
	}

	if user == nil {
		s.logger.Info("User not found")

		return userDomain.ErrUserNotFound
	}

	return nil
}
