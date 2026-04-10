package authDomain

import (
	userDomain "app/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string          `json:"user_id"`
	Role   userDomain.Role `json:"role"`
	jwt.RegisteredClaims
}
