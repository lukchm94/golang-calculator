package authService

import (
	"app/cmd/config"
	authDomain "app/internal/domain/auth"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAuthService struct {
	logger *slog.Logger
	config *config.JwtConfig
}

func NewJwtAuthService(logger *slog.Logger, config *config.JwtConfig) *JwtAuthService {
	return &JwtAuthService{
		logger: logger,
		config: config,
	}
}

func (s *JwtAuthService) GenerateToken(input JwtLoginInput) (string, error) {
	s.logger.Info("Generating JWT token", "userID", input.UserID, "role", input.Role)

	claims := &authDomain.CustomClaims{
		UserID:           input.UserID,
		Role:             input.Role,
		RegisteredClaims: s.buildRegisteredClaims(input),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.SecretKey))

	if err != nil {
		s.logger.Error("Failed to sign JWT token", "error", err)
		return "", err
	}

	return signedToken, nil
}

func (s *JwtAuthService) ValidateToken(tokenStr string) (*authDomain.CustomClaims, error) {
	s.logger.Info("Validating JWT token")

	token, err := jwt.ParseWithClaims(tokenStr, &authDomain.CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Warn("Unexpected signing method", "method", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		s.logger.Error("Failed to parse JWT token", "error", err)
		if err == jwt.ErrSignatureInvalid {
			return nil, authDomain.ErrTokenInvalid
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*authDomain.CustomClaims); ok {
		if !token.Valid {
			s.logger.Warn("Invalid JWT token")
			return nil, authDomain.ErrTokenInvalid
		}

		// Check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			s.logger.Info("JWT token has expired", "userID", claims.UserID)
			return nil, authDomain.ErrTokenExpired
		}

		s.logger.Info("JWT token is valid", "userID", claims.UserID, "role", claims.Role)
		return claims, nil
	}

	s.logger.Warn("Invalid JWT token claims")
	return nil, authDomain.ErrTokenInvalid
}

func (s *JwtAuthService) buildRegisteredClaims(input JwtLoginInput) jwt.RegisteredClaims {
	now := time.Now()
	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(s.config.ExpirationTime)),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    s.config.Issuer.String(),
		Subject:   input.UserID,
	}
}
