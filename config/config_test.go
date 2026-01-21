package config

import (
	"testing"
)

func TestConfig_Struct(t *testing.T) {
	c := Config{
		Port:        "8080",
		Environment: "development",
		SecretKey:   "secret",
		CookieKey:   "cookie-secret",
		Database: DatabaseOptions{
			Dsn:          "postgres://localhost/test",
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  "15m",
		},
		Limiter: LimiterOptions{
			RPS:     2,
			Burst:   4,
			Enabled: true,
		},
		Email: Email{
			Host:     "smtp.example.com",
			Username: "user@example.com",
			Password: "password",
		},
	}

	if c.Port != "8080" {
		t.Errorf("Port = %v, want 8080", c.Port)
	}
	if c.Environment != "development" {
		t.Errorf("Environment = %v, want development", c.Environment)
	}
	if c.Database.MaxOpenConns != 25 {
		t.Errorf("Database.MaxOpenConns = %v, want 25", c.Database.MaxOpenConns)
	}
	if c.Limiter.RPS != 2 {
		t.Errorf("Limiter.RPS = %v, want 2", c.Limiter.RPS)
	}
}

func TestDatabaseOptions(t *testing.T) {
	opts := DatabaseOptions{
		Dsn:          "postgres://user:pass@localhost:5432/db",
		MaxOpenConns: 10,
		MaxIdleConns: 5,
		MaxIdleTime:  "10m",
	}

	if opts.Dsn != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("Dsn = %v, want postgres://user:pass@localhost:5432/db", opts.Dsn)
	}
	if opts.MaxOpenConns != 10 {
		t.Errorf("MaxOpenConns = %v, want 10", opts.MaxOpenConns)
	}
	if opts.MaxIdleConns != 5 {
		t.Errorf("MaxIdleConns = %v, want 5", opts.MaxIdleConns)
	}
	if opts.MaxIdleTime != "10m" {
		t.Errorf("MaxIdleTime = %v, want 10m", opts.MaxIdleTime)
	}
}

func TestLimiterOptions(t *testing.T) {
	opts := LimiterOptions{
		RPS:     5.0,
		Burst:   10,
		Enabled: true,
	}

	if opts.RPS != 5.0 {
		t.Errorf("RPS = %v, want 5.0", opts.RPS)
	}
	if opts.Burst != 10 {
		t.Errorf("Burst = %v, want 10", opts.Burst)
	}
	if !opts.Enabled {
		t.Error("Enabled should be true")
	}
}

func TestLimiterOptions_Disabled(t *testing.T) {
	opts := LimiterOptions{
		RPS:     0,
		Burst:   0,
		Enabled: false,
	}

	if opts.Enabled {
		t.Error("Enabled should be false")
	}
}

func TestEmail(t *testing.T) {
	email := Email{
		Host:     "mail.example.com",
		Username: "sender@example.com",
		Password: "secret123",
	}

	if email.Host != "mail.example.com" {
		t.Errorf("Host = %v, want mail.example.com", email.Host)
	}
	if email.Username != "sender@example.com" {
		t.Errorf("Username = %v, want sender@example.com", email.Username)
	}
	if email.Password != "secret123" {
		t.Errorf("Password = %v, want secret123", email.Password)
	}
}

func TestConfig_ZeroValues(t *testing.T) {
	c := Config{}

	if c.Port != "" {
		t.Errorf("Port should be empty, got %v", c.Port)
	}
	if c.Environment != "" {
		t.Errorf("Environment should be empty, got %v", c.Environment)
	}
	if c.Database.MaxOpenConns != 0 {
		t.Errorf("Database.MaxOpenConns should be 0, got %v", c.Database.MaxOpenConns)
	}
	if c.Limiter.RPS != 0 {
		t.Errorf("Limiter.RPS should be 0, got %v", c.Limiter.RPS)
	}
}

func TestDatabaseOptions_ZeroValues(t *testing.T) {
	opts := DatabaseOptions{}

	if opts.Dsn != "" {
		t.Errorf("Dsn should be empty, got %v", opts.Dsn)
	}
	if opts.MaxOpenConns != 0 {
		t.Errorf("MaxOpenConns should be 0, got %v", opts.MaxOpenConns)
	}
	if opts.MaxIdleConns != 0 {
		t.Errorf("MaxIdleConns should be 0, got %v", opts.MaxIdleConns)
	}
	if opts.MaxIdleTime != "" {
		t.Errorf("MaxIdleTime should be empty, got %v", opts.MaxIdleTime)
	}
}

func TestEmail_ZeroValues(t *testing.T) {
	email := Email{}

	if email.Host != "" {
		t.Errorf("Host should be empty, got %v", email.Host)
	}
	if email.Username != "" {
		t.Errorf("Username should be empty, got %v", email.Username)
	}
	if email.Password != "" {
		t.Errorf("Password should be empty, got %v", email.Password)
	}
}

func TestConfig_NestedStructs(t *testing.T) {
	c := Config{
		Database: DatabaseOptions{
			Dsn:          "postgres://localhost/testdb",
			MaxOpenConns: 50,
			MaxIdleConns: 25,
			MaxIdleTime:  "30m",
		},
		Limiter: LimiterOptions{
			RPS:     10.0,
			Burst:   20,
			Enabled: true,
		},
		Email: Email{
			Host:     "smtp.test.com",
			Username: "test@test.com",
			Password: "testpass",
		},
	}

	// Verify Database
	if c.Database.Dsn != "postgres://localhost/testdb" {
		t.Errorf("Database.Dsn = %v, want postgres://localhost/testdb", c.Database.Dsn)
	}
	if c.Database.MaxOpenConns != 50 {
		t.Errorf("Database.MaxOpenConns = %v, want 50", c.Database.MaxOpenConns)
	}

	// Verify Limiter
	if c.Limiter.RPS != 10.0 {
		t.Errorf("Limiter.RPS = %v, want 10.0", c.Limiter.RPS)
	}
	if c.Limiter.Burst != 20 {
		t.Errorf("Limiter.Burst = %v, want 20", c.Limiter.Burst)
	}

	// Verify Email
	if c.Email.Host != "smtp.test.com" {
		t.Errorf("Email.Host = %v, want smtp.test.com", c.Email.Host)
	}
}

func TestConfig_ProductionSettings(t *testing.T) {
	c := Config{
		Port:        "443",
		Environment: "production",
		SecretKey:   "super-secret-production-key",
		CookieKey:   "32-byte-encryption-key-here!!!!",
		Database: DatabaseOptions{
			Dsn:          "postgres://prod-user:prod-pass@prod-host:5432/proddb?sslmode=require",
			MaxOpenConns: 100,
			MaxIdleConns: 50,
			MaxIdleTime:  "5m",
		},
		Limiter: LimiterOptions{
			RPS:     100,
			Burst:   200,
			Enabled: true,
		},
	}

	if c.Environment != "production" {
		t.Errorf("Environment = %v, want production", c.Environment)
	}
	if c.Database.MaxOpenConns != 100 {
		t.Errorf("Database.MaxOpenConns = %v, want 100", c.Database.MaxOpenConns)
	}
}

func TestConfig_DevelopmentSettings(t *testing.T) {
	c := Config{
		Port:        "8080",
		Environment: "development",
		SecretKey:   "dev-secret",
		CookieKey:   "dev-cookie-key",
		Database: DatabaseOptions{
			Dsn:          "postgres://localhost/devdb?sslmode=disable",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
			MaxIdleTime:  "15m",
		},
		Limiter: LimiterOptions{
			RPS:     1000,
			Burst:   2000,
			Enabled: false, // Disable in dev for easier testing
		},
	}

	if c.Environment != "development" {
		t.Errorf("Environment = %v, want development", c.Environment)
	}
	if c.Limiter.Enabled {
		t.Error("Limiter should be disabled in development")
	}
}
