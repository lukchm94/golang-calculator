# Create Admin Health Check with JWT Authentication

Implement three features:
1. Add Role column to PostgreSQL user table (Admin/User enum)
2. Generate JWT tokens for admin users
3. Create admin health endpoint with JWT validation

---

## Phase 1: Add Role Column to User

### 1.1 Domain Layer
- [ ] Add `Role` field (Enum type: Admin/Guest) to the `User` struct in `internal/domain/user/model.go`
- [ ] Create a Role type definition or enum (e.g., `UserRole` type with `Admin` and `Guest` constants)
- [ ] Add a validation method to ensure only valid role values are used

### 1.2 Persistence Layer
- [ ] Update `UserPostgres` struct in `internal/infrastructure/postgres/models/user.go` to include the `Role` field with GORM PostgreSQL enum type
- [ ] Update the `ToDomain()` method to include the Role field mapping
- [ ] Update the `FromDomain()` method to include the Role field mapping
- [ ] Create a database migration file (e.g., `001_add_role_to_users.sql` or use GORM's migration approach)
- [ ] Update `InitPostgresTables()` in `internal/infrastructure/postgres/client.go` to handle the migration

### 1.3 User Repositories
- [ ] Update user repository methods (if any) to handle the new Role field during create/read operations

---

## Phase 2: JWT Token Generation for Admin Users

### 2.1 Dependencies & Configuration
- [ ] Add JWT library to `go.mod` (e.g., `github.com/golang-jwt/jwt/v5`)
- [ ] Add JWT secret key to environment configuration in `cmd/config/env_config.go`
- [ ] Create a JWT configuration struct to hold secret key, token expiration, issuer, etc.

### 2.2 JWT Service Layer
- [ ] Create a new JWT service package (e.g., `internal/application/auth/jwt_service.go`)
- [ ] Implement JWT token generation method with claims (user ID, email, role, etc.)
- [ ] Implement JWT token validation method
- [ ] Implement JWT token parsing and claims extraction method
- [ ] Add error handling for JWT operations (invalid token, expired token, etc.)

### 2.3 User Service Enhancement
- [ ] Update `internal/application/user/user_service.go` or create auth service to include login functionality
- [ ] Add method to generate JWT token after successful login for admin users
- [ ] Ensure role validation during token generation (only Admin users can generate admin tokens)

### 2.4 User Controller Update
- [ ] Enhance the user controller in `internal/infrastructure/http/controllers/user_controller.go` to handle token generation response

---

## Phase 3: Health Admin Endpoint with JWT Validation

### 3.1 Middleware Layer
- [ ] Create an auth middleware package (e.g., `internal/infrastructure/http/middleware/auth_middleware.go`)
- [ ] Implement JWT validation middleware that extracts and validates tokens from Authorization headers
- [ ] Implement role-based authorization middleware to check if user is Admin
- [ ] Add error responses for missing/invalid/expired tokens

### 3.2 Routes & Endpoints
- [ ] Add new route constant for admin health endpoint in `internal/infrastructure/http/app_routes.go` (e.g., `/health/admin`)
- [ ] Update the router in `internal/infrastructure/http/router.go` to support middleware chain

### 3.3 Health Handler Enhancement
- [ ] Create a new handler method or new handler for admin health endpoint in `internal/infrastructure/http/health_handler.go`
- [ ] Apply JWT middleware and Admin role middleware to the admin endpoint
- [ ] Ensure proper error responses (401 Unauthorized, 403 Forbidden, etc.)

### 3.4 Health Service (if needed)
- [ ] Check if `internal/application/health_service.go` needs enhancement for admin-specific health data

---

## Phase 4: Integration & Setup

### 4.1 Initialization
- [ ] Update server initialization in `cmd/server/main.go` to:
  - Initialize JWT service with config
  - Pass JWT service to user service/controller
  - Wire middleware into router

### 4.2 Error Handling
- [ ] Add JWT-related error types in `internal/domain/user/errors.go`
- [ ] Create HTTP error responses in `internal/infrastructure/http/errors/` for auth-related errors

### 4.3 Configuration
- [ ] Add JWT_SECRET to environment variables (document in `.env` or config file)
- [ ] Add JWT token expiration time configuration
- [ ] Update `cmd/config/env_config.go` to read JWT configuration from environment

---

## Phase 5: Testing & Validation

- [ ] Add role initialization (default to "User") for existing user records or new user registration
- [ ] Test JWT token generation for admin users
- [ ] Test JWT token validation in middleware
- [ ] Test that non-admin users cannot generate admin tokens
- [ ] Test admin health endpoint with valid/invalid/expired tokens
- [ ] Test role-based access control

---

## Key Decisions to Make Before Implementation

1. Should default user registration have `User` role or will admins be manually created?
2. What should the JWT token exp time be? (e.g., 24h, 7d)
3. Should JWT validation error responses include detail or be generic for security?
4. Will you use a separate endpoint for admin token generation or include it in login?
