package domain

type HealthStatusResponse string

const (
	HealthyStatus HealthStatusResponse = "healthy"
	ErrorStatus   HealthStatusResponse = "error"
)

type HealthStatus struct {
	Status HealthStatusResponse `json:"status"`
}

func HealthCheck() HealthStatus {
	return HealthStatus{Status: HealthyStatus}
}
