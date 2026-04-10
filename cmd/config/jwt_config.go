package config

import "time"

type JwtConfig struct {
	SecretKey      string
	Issuer         Issuer
	ExpirationTime time.Duration
}

func FromStringToTimeDuration(expirationTime string) time.Duration {
	d, err := time.ParseDuration(expirationTime)
	if err != nil {
		return 3600 * time.Second
	}
	return d
}
