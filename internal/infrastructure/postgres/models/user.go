package postgresModels

import (
	userDomain "app/internal/domain/user"
	"time"
)

type UserPostgres struct {
	ID             string          `gorm:"primaryKey;type:uuid"`
	FirstName      string          `gorm:"not null"`
	LastName       string          `gorm:"not null"`
	Email          string          `gorm:"uniqueIndex;not null"`
	Role           userDomain.Role `gorm:"type:varchar(20);not null;default:'guest'"`
	HashedPassword string          `gorm:"not null"`
	CreatedAt      time.Time       `gorm:"autoCreateTime"`
}

// TableName overrides the table name for GORM
func (UserPostgres) TableName() string {
	return "users"
}

// ToDomain converts DB model to Domain model
func (up *UserPostgres) ToDomain() *userDomain.User {
	return &userDomain.User{
		ID:             up.ID,
		FirstName:      up.FirstName,
		LastName:       up.LastName,
		Email:          up.Email,
		Role:           up.Role,
		HashedPassword: up.HashedPassword,
		CreatedAt:      up.CreatedAt,
	}
}

// FromDomain converts Domain model to DB model
func FromDomain(u *userDomain.User) *UserPostgres {
	return &UserPostgres{
		ID:             u.ID,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Email:          u.Email,
		Role:           u.Role,
		HashedPassword: u.HashedPassword,
		CreatedAt:      u.CreatedAt,
	}
}
