package authService

import userDomain "app/internal/domain/user"

type JwtLoginInput struct {
	UserID string          `json:"userID" validate:"required"`
	Role   userDomain.Role `json:"role" validate:"required"`
}
