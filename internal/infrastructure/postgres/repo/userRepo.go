package postgresRepo

import (
	"errors"
	"log/slog"

	userDomain "app/internal/domain/user"
	postgresModels "app/internal/infrastructure/postgres/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

	if err := r.db.Create(postgresUser).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			r.logger.Info("User with email already exists", "email", user.Email)
			return userDomain.ErrUserAlreadyExists
		}

		r.logger.Error("Failed to create user", "error", err)
		return err
	}

	return nil
}

func (r *UserRepository) GenerateUserID() string {
	return uuid.New().String()
}
