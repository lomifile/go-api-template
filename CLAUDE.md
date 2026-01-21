# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Run the application
go run ./cmd/app/main.go

# Build binary
go build -o api ./cmd/app/main.go

# Lint (uses golangci-lint with custom config)
golangci-lint run

# Format code
golangci-lint fmt

# Tidy dependencies
go mod tidy
```

No tests exist yet. When adding tests, use standard Go testing conventions with `go test ./...`.

## Architecture

This is a Go REST API server using Fiber v2, following clean architecture principles.

### Layer Structure

```
cmd/app/           → Entry point and application bootstrap
api/http/          → HTTP layer (handlers, middleware, router)
internal/domain/   → Business logic (models, services, repositories)
internal/adapter/  → Database adapters
pkg/               → Shared utilities (logger, postgres, email, utils)
config/            → Configuration parsing
```

### Request Flow

1. **main.go** creates config, logger, and database connection
2. **app.go** configures Fiber with middleware stack:
   - Logger → RequestID → Recover → Helmet → Limiter → CORS → EncryptCookie
3. Routes dispatch to handlers in `api/http/handler/`
4. Handlers call services in `internal/domain/service/`
5. Services use repositories in `internal/domain/repository/`
6. Repositories use the PostgreSQL adapter in `internal/adapter/`

### Key Patterns

- **Dependency Injection**: Config and logger passed through function parameters
- **Option Pattern**: Used in postgres connection (`pkg/postgres/`) and server config
- **Adapter Pattern**: PostgreSQL adapter wraps pgxpool with sqlx support
- **Custom JSON Marshaling**: Domain models implement `sql.Scanner` and `driver.Valuer` for PostgreSQL JSON columns

### Configuration

Environment variables loaded from `.env` file:
- `PORT` - Server port
- `ENVIRONMENT` - development/production
- `DB_URL` - PostgreSQL connection string
- `EMAIL_SERVER`, `EMAIL_USERNAME`, `EMAIL_PASSWORD` - SMTP config
- `COOKIE_KEY` - Cookie encryption key

All config values support flag overrides.

### Linting Rules

Uses golangci-lint with:
- Max line length: 100 characters
- Enabled linters: govet, errcheck, staticcheck, ineffassign, revive
- Formatters: gofmt, goimports, golines

### Response Utilities

Use utilities from `pkg/utils/response.go`:
- `SuccessResponseMap[T]` - Generic success wrapper
- `ErrorResponseMap` - Error with request ID
- `PaginationResponse[T]` - Paginated responses
- `ErrorResponder` - Centralized error handler with logging
