package config

type ValidEnvironments string

const (
	DevEnvironment   ValidEnvironments = "dev"
	StageEnvironment ValidEnvironments = "stage"
	ProdEnvironment  ValidEnvironments = "prod"
)

func (e ValidEnvironments) IsValid() bool {
	switch e {
	case DevEnvironment, StageEnvironment, ProdEnvironment:
		return true
	default:
		return false
	}
}

func (e ValidEnvironments) String() string {
	return string(e)
}

func FromString(env string) ValidEnvironments {
	switch env {
	case DevEnvironment.String():
		return DevEnvironment
	case StageEnvironment.String():
		return StageEnvironment
	case ProdEnvironment.String():
		return ProdEnvironment
	default:
		return DevEnvironment
	}
}
