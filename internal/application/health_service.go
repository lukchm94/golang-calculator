package application

import domain "app/internal/domain/health"

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

// This function `CheckHealth` is a method of the `HealthService` struct. It returns a value of type
// `domain.HealthStatus`. Inside the function, it calls the `HealthCheck` function from the `domain`
// package to perform a health check and return the result.
func (s *HealthService) CheckHealth() domain.HealthStatus {
	healthStatus := domain.HealthCheck()

	println("Health check performed, status:", healthStatus.Status)
	return healthStatus
}
