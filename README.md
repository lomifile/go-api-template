# Go Server Template

A production-ready Go REST API server template using Fiber v2 framework with clean architecture principles. Use this as a starting point for building your own REST APIs.

## Prerequisites

- Go 1.25.1 or higher
- PostgreSQL database
- [golangci-lint](https://golangci-lint.run/usage/install/) (for linting)

## Getting Started

### 1. Clone and configure

```bash
cp .env.example .env  # or edit .env directly
```

### 2. Configure environment variables

Edit `.env` with your settings:

```env
ENVIRONMENT=development
PORT=8080
DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
EMAIL_SERVER=smtp.example.com
EMAIL_USERNAME=your-email@example.com
EMAIL_PASSWORD=your-password
COOKIE_KEY=32-byte-secret-key-for-cookies
```

### 3. Run the application

```bash
go run ./cmd/app/main.go
```

The server starts on the configured port (default: 8080).

## Development Commands

```bash
# Run application
go run ./cmd/app/main.go

# Build binary
go build -o api ./cmd/app/main.go

# Run linter
golangci-lint run

# Format code
golangci-lint fmt

# Tidy dependencies
go mod tidy

# Run tests (when added)
go test ./...
```

## Configuration

All configuration options can be set via environment variables or command-line flags.

### Environment Variables

| Variable         | Description                      | Default |
| ---------------- | -------------------------------- | ------- |
| `PORT`           | Server port                      | `8080`  |
| `ENVIRONMENT`    | `development` or `production`    | -       |
| `DB_URL`         | PostgreSQL connection string     | -       |
| `SECRET_KEY`     | JWT secret key                   | -       |
| `COOKIE_KEY`     | Cookie encryption key (32 bytes) | -       |
| `EMAIL_SERVER`   | SMTP server hostname             | -       |
| `EMAIL_USERNAME` | SMTP username                    | -       |
| `EMAIL_PASSWORD` | SMTP password                    | -       |

### Command-Line Flags

```bash
go run ./cmd/app/main.go \
  -port=8080 \
  -environment=development \
  -db-dsn="postgres://..." \
  -db-max-open-conns=25 \
  -db-max-idle-conns=25 \
  -db-max-idle-time=15m \
  -limiter-rps=2 \
  -limiter-burst=4 \
  -limiter-enabled=true
```

## Project Structure

```
.
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── api/
│   └── http/
│       ├── handler/             # HTTP request handlers
│       ├── middleware/          # Custom middleware (logger, etc.)
│       └── router/              # Route definitions
├── config/
│   └── config.go                # Configuration management
├── internal/
│   ├── app/
│   │   └── app.go               # Application bootstrap and middleware setup
│   ├── server/
│   │   └── server.go            # Fiber server wrapper with graceful shutdown
│   ├── adapter/
│   │   └── postgres.go          # Database adapter (pgxpool + sqlx)
│   └── domain/
│       ├── model/               # Domain models
│       ├── repository/          # Repository interfaces
│       └── service/             # Business logic services
├── pkg/
│   ├── logger/                  # Zap logger wrapper
│   ├── postgres/                # PostgreSQL connection pool with retry
│   ├── email/                   # SMTP email client
│   └── utils/                   # Response types, pagination, helpers
├── static/                      # Static assets
├── .env                         # Environment configuration
├── .golangci.yml                # Linter configuration
└── go.mod                       # Go module definition
```

## Architecture

This template follows clean architecture principles with clear separation of concerns:

### Layers

1. **HTTP Layer** (`api/http/`) - Handles HTTP requests, routing, and middleware
2. **Domain Layer** (`internal/domain/`) - Business logic, models, and repository interfaces
3. **Infrastructure Layer** (`internal/adapter/`, `pkg/`) - Database adapters, external services

### Request Flow

```
HTTP Request
    ↓
Middleware Stack (Logger → RequestID → Recover → Helmet → Limiter → CORS → EncryptCookie)
    ↓
Router
    ↓
Handler
    ↓
Service (business logic)
    ↓
Repository (data access)
    ↓
Database Adapter
    ↓
PostgreSQL
```

### Middleware Stack

The application uses the following middleware in order:

1. **Logger** - Structured request logging with request ID, latency, status
2. **RequestID** - Generates unique request identifiers
3. **Recover** - Panic recovery
4. **Helmet** - Security headers
5. **Limiter** - Rate limiting (100 requests/minute)
6. **CORS** - Cross-origin resource sharing
7. **EncryptCookie** - Cookie encryption

## Server Configuration

Default server settings:

| Setting          | Value               |
| ---------------- | ------------------- |
| Read Timeout     | 10 seconds          |
| Write Timeout    | 5 seconds           |
| Shutdown Timeout | 3 seconds           |
| Body Limit       | 32 MB               |
| Rate Limit       | 100 requests/minute |

## Database

### Connection Pool

The PostgreSQL connection pool includes:

- Automatic retry logic (10 attempts by default)
- Connection timeout: 1 second
- Max connections: 4 (configurable)
- DSN validation (must use `postgres://` or `postgresql://` scheme)

### Adapter

The database adapter combines `pgxpool` for connection pooling with `sqlx` for convenient query scanning:

```go
db, err := adapter.NewPostgresAdapter(postgresPool)
// Use db.Select, db.Get, db.Exec, etc.
```

## Logging

Uses [Zap](https://github.com/uber-go/zap) for structured logging:

- **Development mode**: Colorized console output
- **Production mode**: JSON format

HTTP request logs include:

- `request_id` - Unique request identifier
- `http_method` - GET, POST, etc.
- `http_path` - Request path
- `status` - Response status code
- `latency_ms` - Request duration
- `client_ip` - Client IP address
- `user_agent` - User agent string

## API Response Format

### Success Response

```json
{
  "request_id": "uuid",
  "status": 200,
  "data": { ... },
  "ts": "2024-01-01T00:00:00Z"
}
```

### Error Response

```json
{
  "request_id": "uuid",
  "status": 400,
  "error": "error message",
  "ts": "2024-01-01T00:00:00Z"
}
```

### Paginated Response

```json
{
  "items": [...],
  "total": 100,
  "meta": {
    "next": 2,
    "previous": null,
    "has_next_page": true,
    "has_previous_page": false
  }
}
```

## Graceful Shutdown

The server handles `SIGTERM` and `SIGINT` signals for graceful shutdown:

1. Stops accepting new connections
2. Waits for active requests to complete (up to 3 seconds)
3. Closes database connections
4. Flushes logger

## Linting

The project uses golangci-lint with the following configuration:

- **Max line length**: 100 characters
- **Formatters**: gofmt, goimports, golines
- **Linters**: govet, errcheck, staticcheck, ineffassign, revive

Run with:

```bash
golangci-lint run
```

## Extending the Template

### Adding a New Endpoint

1. Define the model in `internal/domain/model/`
2. Create repository interface in `internal/domain/repository/`
3. Implement service in `internal/domain/service/`
4. Create handler in `api/http/handler/`
5. Register route in `api/http/router/router.go`

### Adding Custom Middleware

Create middleware in `api/http/middleware/` following the Fiber handler pattern:

```go
func MyMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Before request
        err := c.Next()
        // After request
        return err
    }
}
```
