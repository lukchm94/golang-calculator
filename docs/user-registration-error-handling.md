## Improve user registration error handling and HTTP responses

### Description

The current user registration flow works end-to-end but has several rough edges around error semantics and HTTP responses. We should make the behavior more consistent and expressive across the controller, application service, and Postgres repository, and return more accurate HTTP status codes to clients.

This affects:

- `user_controller`
- `user_handler`
- `userService`
- `userDomain` errors
- Postgres `UserRepository` implementation

---

### Tasks

- **Controller / HTTP layer**
  - **Method handling**
    - Update `UserController.validateRegisterReq` to return a real error when `r.Method != http.MethodPost` instead of `nil`.
    - Introduce a specific error type (for example `MethodNotAllowedError`) in `internal/infrastructure/http/errors`.
    - Update `UserHandler.handleErrors` to map that error to `405 Method Not Allowed`.
  - **Validation responses**
    - Ensure `InvalidRequestError` and `MissingFieldError` are mapped to `400 Bad Request` (or `422 Unprocessable Entity`) instead of always returning `500`.
    - Return a structured JSON error response (for example `{"error":"...", "field":"..."}` for `MissingFieldError`).
    - Confirm that invalid JSON (decode errors) and missing required JSON fields produce clear messages.

- **User service (`UserService.Register`)**
  - **Duplicate email handling**
    - Define a domain-level error `ErrEmailAlreadyInUse` in `internal/domain/user`.
    - Ensure the Postgres repo maps the DB unique-constraint violation on `email` to `ErrEmailAlreadyInUse`.
    - Update `UserHandler.handleErrors` to translate `ErrEmailAlreadyInUse` into `409 Conflict`.

- **Postgres user repository**
  - **Not-found semantics**
    - Decide and implement a consistent contract for `GetUserByID` and `GetUserByEmail`:
      - Option A: return `nil, nil` for “not found” and let the service map that to `ErrUserNotFound`.
      - Option B: keep returning `ErrUserNotFound`, and adjust the service to treat that as an expected, non-500 condition.
    - Ensure logging does not log expected “not found” conditions as `ERROR`.
  - **Error mapping**
    - Map GORM `ErrRecordNotFound` (and duplicate-key errors) to clear domain errors instead of leaking raw DB errors.

- **HTTP error mapping (centralized)**
  - Enhance `UserHandler.handleErrors` to:
    - Pattern-match on error types (`InvalidRequestError`, `MissingFieldError`, `MethodNotAllowedError`, `ErrEmailAlreadyInUse`, `ErrUserNotFound`, generic errors).
    - Set appropriate `StatusCode` values:
      - `400` / `422` for invalid or missing data.
      - `405` for unsupported method.
      - `409` for email already in use.
      - `500` as a fallback for unexpected errors.
    - Return a consistent error JSON shape for all error responses.

---

### Acceptance criteria

- Registering a user with a **valid JSON body** returns:
  - `201 Created`
  - JSON with the new user’s ID and basic fields.
- Sending:
  - **Non-JSON** or malformed JSON → `400`/`422` with a clear `"error": "invalid request"` style message.
  - Missing `first_name` / `last_name` / `email` / `password` → `400`/`422` with `"error": "missing required field"` and the field name included.
  - Duplicate `email` (violates unique constraint) → `409 Conflict` with a clear error like `"email already in use"`.
  - Wrong HTTP method on `/register` (for example `GET`) → `405 Method Not Allowed`.
- Logs:
  - Do not log expected validation or “not found” situations as `ERROR`; reserve `ERROR` for real server or DB issues.

