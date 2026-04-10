package domain

import userDomain "app/internal/domain/user"

type HealthStatusResponse string

const (
	HealthyStatus HealthStatusResponse = "healthy"
	ErrorStatus   HealthStatusResponse = "error"
)

type HealthStatus struct {
	Status HealthStatusResponse `json:"status"`
}

type AdminHealthStatus struct {
	Status   HealthStatusResponse `json:"status"`
	Username string               `json:"username"`
	Email    string               `json:"email"`
	Role     userDomain.Role
}

func HealthCheck() HealthStatus {
	return HealthStatus{Status: HealthyStatus}
}

func AdminHealthCheck(user *userDomain.User) AdminHealthStatus {
	return AdminHealthStatus{
		Status:   HealthyStatus,
		Username: user.FullName(),
		Email:    user.Email,
		Role:     user.Role,
	}
}
