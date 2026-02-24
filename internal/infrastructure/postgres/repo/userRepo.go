package postgresRepo

import (
	"log/slog"

	userDomain "app/internal/domain/users"
	postgresModels "app/internal/infrastructure/postgres/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) GetUserByID(id string) (*userDomain.User, error) {
	var user postgresModels.UserPostgres
	r.logger.Debug("Getting user by ID", "id", id)

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Info("User not found", "id", id)
			return nil, userDomain.ErrUserNotFound
		}

		r.logger.Error("Failed to get user by ID", "error", err)
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *UserRepository) GetUserByEmail(email string) (*userDomain.User, error) {
	var user postgresModels.UserPostgres

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Info("User not found", "email", email)
			return nil, userDomain.ErrUserNotFound
		}

		r.logger.Error("Failed to get user by email", "error", err)
		return nil, err
	}

	r.logger.Info("User found", "email", email)
	return user.ToDomain(), nil
}

func (r *UserRepository) CreateUser(user *userDomain.User) error {
	postgresUser := postgresModels.FromDomain(user)

	r.logger.Debug("Creating user", "user", postgresUser)

	return r.db.Create(postgresUser).Error
}

func (r *UserRepository) GenerateUserID() string {
	return uuid.New().String()
}
