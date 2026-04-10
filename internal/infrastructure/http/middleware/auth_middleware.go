package middleware

import (
	authService "app/internal/application/auth"
	authDomain "app/internal/domain/auth"
	userDomain "app/internal/domain/user"
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

func AuthMiddleware(logger *slog.Logger, jwtService *authService.JwtAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing authorization header")
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("Invalid authorization header format")
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				logger.Warn("Token validation failed", "error", err)
				if err == authDomain.ErrTokenExpired {
					http.Error(w, "Token has expired", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Store claims in context
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnlyMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserClaimsKey).(*authDomain.CustomClaims)
			if !ok {
				logger.Error("Missing user claims in context")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if userDomain.Role(claims.Role) != userDomain.Admin {
				logger.Warn("Non-admin user attempted admin access", "userID", claims.UserID)
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
