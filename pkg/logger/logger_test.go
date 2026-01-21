package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestNew_Development(t *testing.T) {
	l := New(Config{
		Debug: true,
		Env:   "development",
	})

	if l == nil {
		t.Fatal("New() returned nil")
	}
	if l.Logger == nil {
		t.Fatal("Logger.Logger is nil")
	}

	// Should not panic
	l.Info("test message")
	_ = l.Sync()
}

func TestNew_Production(t *testing.T) {
	l := New(Config{
		Debug: false,
		Env:   "production",
	})

	if l == nil {
		t.Fatal("New() returned nil")
	}

	// Should not panic
	l.Info("test message")
	_ = l.Sync()
}

func TestNew_WithService(t *testing.T) {
	l := New(Config{
		Debug:   true,
		Service: "test-service",
		Env:     "test",
	})

	if l == nil {
		t.Fatal("New() returned nil")
	}

	// Should not panic
	l.Info("test message with service")
	_ = l.Sync()
}

func TestLogger_Sugar(t *testing.T) {
	l := New(Config{Debug: true})

	sugar := l.Sugar()
	if sugar == nil {
		t.Fatal("Sugar() returned nil")
	}

	// Should not panic
	sugar.Info("sugared log message")
	_ = l.Sync()
}

func TestLogger_Named(t *testing.T) {
	l := New(Config{Debug: true})

	named := l.Named("sub-component")
	if named == nil {
		t.Fatal("Named() returned nil")
	}
	if named.Logger == nil {
		t.Fatal("Named logger's internal Logger is nil")
	}

	// Should not panic
	named.Info("named log message")
	_ = l.Sync()
}

func TestLogger_With(t *testing.T) {
	l := New(Config{Debug: true})

	withFields := l.With(
		zap.String("key1", "value1"),
		zap.Int("key2", 42),
	)

	if withFields == nil {
		t.Fatal("With() returned nil")
	}
	if withFields.Logger == nil {
		t.Fatal("With logger's internal Logger is nil")
	}

	// Should not panic
	withFields.Info("log with fields")
	_ = l.Sync()
}

func TestLogger_Chaining(t *testing.T) {
	l := New(Config{Debug: true, Service: "chain-test"})

	chained := l.Named("component").With(zap.String("request_id", "123"))

	if chained == nil {
		t.Fatal("Chained logger is nil")
	}

	// Should not panic
	chained.Info("chained log message")
	_ = l.Sync()
}

func TestLogger_LogLevels(t *testing.T) {
	l := New(Config{Debug: true})

	// These should not panic
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")

	_ = l.Sync()
}

func TestConfig_Empty(t *testing.T) {
	l := New(Config{})

	if l == nil {
		t.Fatal("New() with empty config returned nil")
	}

	l.Info("empty config test")
	_ = l.Sync()
}
