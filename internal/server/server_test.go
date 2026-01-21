package server

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestNew_Default(t *testing.T) {
	s := New()

	if s == nil {
		t.Fatal("New() returned nil")
	}
	if s.App == nil {
		t.Fatal("Server.App is nil")
	}
	if s.address != _defaultAddr {
		t.Errorf("address = %v, want %v", s.address, _defaultAddr)
	}
	if s.readTimeout != _defaultReadTimeout {
		t.Errorf("readTimeout = %v, want %v", s.readTimeout, _defaultReadTimeout)
	}
	if s.writeTimeout != _defaultWriteTimeout {
		t.Errorf("writeTimeout = %v, want %v", s.writeTimeout, _defaultWriteTimeout)
	}
	if s.shutdownTimeout != _defaultShutdownTimeout {
		t.Errorf("shutdownTimeout = %v, want %v", s.shutdownTimeout, _defaultShutdownTimeout)
	}
}

func TestNew_WithPort(t *testing.T) {
	s := New(Port("8080"))

	if s.address != ":8080" {
		t.Errorf("address = %v, want :8080", s.address)
	}
}

func TestPort_Option(t *testing.T) {
	tests := []struct {
		port     string
		expected string
	}{
		{"8080", ":8080"},
		{"3000", ":3000"},
		{"443", ":443"},
	}

	for _, tt := range tests {
		t.Run(tt.port, func(t *testing.T) {
			s := &Server{}
			opt := Port(tt.port)
			opt(s)

			if s.address != tt.expected {
				t.Errorf("address = %v, want %v", s.address, tt.expected)
			}
		})
	}
}

func TestServer_FiberConfig(t *testing.T) {
	s := New()

	// Test that Fiber app is properly configured
	if s.App == nil {
		t.Fatal("Fiber app is nil")
	}

	// Add a test route
	s.App.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Test the route
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := s.App.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Response = %v, want OK", string(body))
	}
}

func TestServer_Notify(t *testing.T) {
	s := New()

	notify := s.Notify()
	if notify == nil {
		t.Error("Notify() returned nil channel")
	}
}

func TestDefaults(t *testing.T) {
	if _defaultAddr != ":80" {
		t.Errorf("_defaultAddr = %v, want :80", _defaultAddr)
	}
	if _defaultReadTimeout != 10*time.Second {
		t.Errorf("_defaultReadTimeout = %v, want %v", _defaultReadTimeout, 10*time.Second)
	}
	if _defaultWriteTimeout != 5*time.Second {
		t.Errorf("_defaultWriteTimeout = %v, want %v", _defaultWriteTimeout, 5*time.Second)
	}
	if _defaultShutdownTimeout != 3*time.Second {
		t.Errorf("_defaultShutdownTimeout = %v, want %v", _defaultShutdownTimeout, 3*time.Second)
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	s := New(Port("0")) // Port 0 lets OS assign available port

	// Start server in background
	s.Start()

	// Give it a moment to start
	time.Sleep(50 * time.Millisecond)

	// Shutdown should not error
	err := s.Shutdown()
	if err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

func TestServer_MultipleOptions(t *testing.T) {
	s := New(
		Port("9090"),
	)

	if s.address != ":9090" {
		t.Errorf("address = %v, want :9090", s.address)
	}
}

func TestServer_RouteHandling(t *testing.T) {
	s := New()

	s.App.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	s.App.Post("/data", func(c *fiber.Ctx) error {
		return c.SendStatus(201)
	})

	// Test GET
	req := httptest.NewRequest("GET", "/health", nil)
	resp, _ := s.App.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("GET /health status = %d, want 200", resp.StatusCode)
	}

	// Test POST
	req = httptest.NewRequest("POST", "/data", nil)
	resp, _ = s.App.Test(req)
	if resp.StatusCode != 201 {
		t.Errorf("POST /data status = %d, want 201", resp.StatusCode)
	}

	// Test 404
	req = httptest.NewRequest("GET", "/notfound", nil)
	resp, _ = s.App.Test(req)
	if resp.StatusCode != 404 {
		t.Errorf("GET /notfound status = %d, want 404", resp.StatusCode)
	}
}
