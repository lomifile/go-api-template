package middleware

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/lomifile/api/pkg/logger"
)

func TestLoggerMiddleware(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(LoggerMiddleware(l))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Status = %d, want 200", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Body = %v, want OK", string(body))
	}

	// Verify request ID header is set
	reqID := resp.Header.Get(fiber.HeaderXRequestID)
	if reqID == "" {
		t.Error("X-Request-ID header should be set")
	}

	_ = l.Sync()
}

func TestLoggerMiddleware_DifferentStatusCodes(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})

	tests := []struct {
		name       string
		path       string
		statusCode int
	}{
		{"200 OK", "/ok", 200},
		{"201 Created", "/created", 201},
		{"400 Bad Request", "/bad", 400},
		{"404 Not Found", "/notfound", 404},
		{"500 Server Error", "/error", 500},
	}

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(LoggerMiddleware(l))

	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Get("/created", func(c *fiber.Ctx) error {
		return c.SendStatus(201)
	})
	app.Get("/bad", func(c *fiber.Ctx) error {
		return c.SendStatus(400)
	})
	app.Get("/notfound", func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
	app.Get("/error", func(c *fiber.Ctx) error {
		return c.SendStatus(500)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("Status = %d, want %d", resp.StatusCode, tt.statusCode)
			}
		})
	}

	_ = l.Sync()
}

func TestLoggerMiddleware_DifferentMethods(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(LoggerMiddleware(l))

	app.Get("/resource", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Post("/resource", func(c *fiber.Ctx) error {
		return c.SendStatus(201)
	})
	app.Put("/resource", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Delete("/resource", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	methods := []struct {
		method     string
		wantStatus int
	}{
		{"GET", 200},
		{"POST", 201},
		{"PUT", 200},
		{"DELETE", 204},
	}

	for _, m := range methods {
		t.Run(m.method, func(t *testing.T) {
			req := httptest.NewRequest(m.method, "/resource", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != m.wantStatus {
				t.Errorf("Status = %d, want %d", resp.StatusCode, m.wantStatus)
			}
		})
	}

	_ = l.Sync()
}

func TestLoggerMiddleware_WithUserAgent(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(LoggerMiddleware(l))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "TestAgent/1.0")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Status = %d, want 200", resp.StatusCode)
	}

	_ = l.Sync()
}

func TestLoggerMiddleware_RequestIDFromHeader(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(LoggerMiddleware(l))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	customReqID := "custom-request-id-123"
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(fiber.HeaderXRequestID, customReqID)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	// Request ID middleware should use the provided ID
	respReqID := resp.Header.Get(fiber.HeaderXRequestID)
	if respReqID != customReqID {
		t.Errorf("X-Request-ID = %v, want %v", respReqID, customReqID)
	}

	_ = l.Sync()
}

func TestLoggerMiddleware_ChainedHandlers(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(LoggerMiddleware(l))

	// Add another middleware after logger
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Custom-Header", "test")
		return c.Next()
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Status = %d, want 200", resp.StatusCode)
	}

	customHeader := resp.Header.Get("X-Custom-Header")
	if customHeader != "test" {
		t.Errorf("X-Custom-Header = %v, want test", customHeader)
	}

	_ = l.Sync()
}
