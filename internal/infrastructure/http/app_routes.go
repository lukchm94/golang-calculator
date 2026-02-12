package httpInfra

type AppRoutes string

const (
	HealthRoute    AppRoutes = "/health"
	CalculateRoute AppRoutes = "/calculator"
)
