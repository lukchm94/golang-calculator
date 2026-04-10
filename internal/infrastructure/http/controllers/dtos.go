package controllers

import (
	userDomain "app/internal/domain/user"
	"time"
)

type UserLoginResponse struct {
	ID        string          `json:"id"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Email     string          `json:"email"`
	Role      userDomain.Role `json:"role"`
	Token     string          `json:"token,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}
