package config

type Issuer string

const (
	AuthServiceIssuer Issuer = "auth-service"
	AppIssuer         Issuer = "golang-calculator"
)

func (i Issuer) IsValid() bool {
	switch i {
	case AuthServiceIssuer, AppIssuer:
		return true
	default:
		return false
	}
}

func (i Issuer) String() string {
	return string(i)
}
