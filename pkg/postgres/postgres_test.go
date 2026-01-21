package postgres

import (
	"testing"
	"time"
)

func TestWithMaxConns(t *testing.T) {
	p := &Postgres{}
	opt := WithMaxConns(10)
	opt(p)

	if p.maxConns != 10 {
		t.Errorf("maxConns = %d, want 10", p.maxConns)
	}
}

func TestWithConnAttempts(t *testing.T) {
	p := &Postgres{}
	opt := WithConnAttempts(5)
	opt(p)

	if p.connAttempts != 5 {
		t.Errorf("connAttempts = %d, want 5", p.connAttempts)
	}
}

func TestWithConnTimeout(t *testing.T) {
	p := &Postgres{}
	timeout := 5 * time.Second
	opt := WithConnTimeout(timeout)
	opt(p)

	if p.connTimeout != timeout {
		t.Errorf("connTimeout = %v, want %v", p.connTimeout, timeout)
	}
}

func TestNew_InvalidDSN(t *testing.T) {
	_, err := New("invalid-dsn")
	if err == nil {
		t.Error("New() with invalid DSN should return error")
	}
}

func TestNew_InvalidScheme(t *testing.T) {
	_, err := New("mysql://localhost:3306/db")
	if err == nil {
		t.Error("New() with invalid scheme should return error")
	}
}

func TestNew_EmptyDSN(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Error("New() with empty DSN should return error")
	}
}

func TestNew_HTTPScheme(t *testing.T) {
	_, err := New("http://localhost:5432/db")
	if err == nil {
		t.Error("New() with http scheme should return error")
	}
}

func TestNew_ValidSchemePostgres(t *testing.T) {
	// This will fail to connect but should pass scheme validation
	_, err := New("postgres://localhost:5432/db", WithConnAttempts(1), WithConnTimeout(100*time.Millisecond))
	if err == nil {
		t.Skip("Skipping: postgres server is available")
	}
	// Error should be about connection, not scheme
	if err.Error() == "postgres: invalid scheme: postgres" {
		t.Error("postgres:// scheme should be valid")
	}
}

func TestNew_ValidSchemePostgresql(t *testing.T) {
	// This will fail to connect but should pass scheme validation
	_, err := New("postgresql://localhost:5432/db", WithConnAttempts(1), WithConnTimeout(100*time.Millisecond))
	if err == nil {
		t.Skip("Skipping: postgres server is available")
	}
	// Error should be about connection, not scheme
	if err.Error() == "postgres: invalid scheme: postgresql" {
		t.Error("postgresql:// scheme should be valid")
	}
}

func TestClose_Nil(t *testing.T) {
	var p *Postgres
	// Should not panic
	p.Close()
}

func TestClose_NilPool(t *testing.T) {
	p := &Postgres{Pool: nil}
	// Should not panic
	p.Close()
}

func TestDefaults(t *testing.T) {
	if _defaultMaxConns != 4 {
		t.Errorf("_defaultMaxConns = %d, want 4", _defaultMaxConns)
	}
	if _defaultConnAttempts != 10 {
		t.Errorf("_defaultConnAttempts = %d, want 10", _defaultConnAttempts)
	}
	if _defaultConnTimeout != time.Second {
		t.Errorf("_defaultConnTimeout = %v, want %v", _defaultConnTimeout, time.Second)
	}
}

func TestMultipleOptions(t *testing.T) {
	p := &Postgres{}

	opts := []Option{
		WithMaxConns(20),
		WithConnAttempts(3),
		WithConnTimeout(2 * time.Second),
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.maxConns != 20 {
		t.Errorf("maxConns = %d, want 20", p.maxConns)
	}
	if p.connAttempts != 3 {
		t.Errorf("connAttempts = %d, want 3", p.connAttempts)
	}
	if p.connTimeout != 2*time.Second {
		t.Errorf("connTimeout = %v, want %v", p.connTimeout, 2*time.Second)
	}
}
