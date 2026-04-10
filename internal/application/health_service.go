package application

import (
	domain "app/internal/domain/health"
	userDomain "app/internal/domain/user"
	"log/slog"
)

type HealthService struct {
	logger *slog.Logger
}

func NewHealthService(logger *slog.Logger) *HealthService {
	return &HealthService{logger: logger}
}

// This function `CheckHealth` is a method of the `HealthService` struct. It returns a value of type
// `domain.HealthStatus`. Inside the function, it calls the `HealthCheck` function from the `domain`
// package to perform a health check and return the result.
func (s *HealthService) CheckHealth() domain.HealthStatus {
	healthStatus := domain.HealthCheck()

	s.logger.Info("Health check performed", "status", healthStatus.Status)
	return healthStatus
}

func (s *HealthService) CheckAdminHealth(user *userDomain.User) domain.AdminHealthStatus {
	healthStatus := domain.AdminHealthCheck(user)

	s.logger.Info("Admin health check performed", "status", healthStatus.Status)
	return healthStatus
}
