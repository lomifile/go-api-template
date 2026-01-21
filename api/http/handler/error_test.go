package handler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/lomifile/api/pkg/logger"
	"github.com/lomifile/api/pkg/utils"
	"go.uber.org/zap"
)

func TestNewErrorResponder(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})

	responder := NewErrorResponder(l)
	if responder == nil {
		t.Fatal("NewErrorResponder() returned nil")
	}
	if responder.l != l {
		t.Error("ErrorResponder.l not set correctly")
	}

	_ = l.Sync()
}

func TestErrorResponder_Error(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})
	responder := NewErrorResponder(l)

	app := fiber.New()
	app.Use(requestid.New())

	app.Get("/error", func(c *fiber.Ctx) error {
		return responder.Error(c, 400, "Bad request", "test_error")
	})

	// Provide a request ID header since the error handler reads from request headers
	customReqID := "test-req-id-123"
	req := httptest.NewRequest("GET", "/error", nil)
	req.Header.Set(fiber.HeaderXRequestID, customReqID)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("Status = %d, want 400", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errorResp utils.ErrorResponseMap
	if err := json.Unmarshal(body, &errorResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if errorResp.Status != 400 {
		t.Errorf("Response.Status = %d, want 400", errorResp.Status)
	}
	if errorResp.Error != "Bad request" {
		t.Errorf("Response.Error = %v, want 'Bad request'", errorResp.Error)
	}
	if errorResp.RequestID != customReqID {
		t.Errorf("Response.RequestID = %v, want %v", errorResp.RequestID, customReqID)
	}
	if errorResp.TS == "" {
		t.Error("Response.TS should not be empty")
	}

	_ = l.Sync()
}

func TestErrorResponder_Error_DifferentStatusCodes(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})
	responder := NewErrorResponder(l)

	tests := []struct {
		name    string
		status  int
		message string
		logKey  string
	}{
		{"Bad Request", 400, "Invalid input", "bad_request"},
		{"Unauthorized", 401, "Authentication required", "unauthorized"},
		{"Forbidden", 403, "Access denied", "forbidden"},
		{"Not Found", 404, "Resource not found", "not_found"},
		{"Internal Error", 500, "Internal server error", "internal_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(requestid.New())

			app.Get("/test", func(c *fiber.Ctx) error {
				return responder.Error(c, tt.status, tt.message, tt.logKey)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.status {
				t.Errorf("Status = %d, want %d", resp.StatusCode, tt.status)
			}

			body, _ := io.ReadAll(resp.Body)
			var errorResp utils.ErrorResponseMap
			if err := json.Unmarshal(body, &errorResp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if errorResp.Error != tt.message {
				t.Errorf("Response.Error = %v, want %v", errorResp.Error, tt.message)
			}
		})
	}

	_ = l.Sync()
}

func TestErrorResponder_Error_WithFields(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})
	responder := NewErrorResponder(l)

	app := fiber.New()
	app.Use(requestid.New())

	app.Get("/error-with-fields", func(c *fiber.Ctx) error {
		return responder.Error(
			c,
			422,
			"Validation failed",
			"validation_error",
			zap.String("field", "email"),
			zap.String("reason", "invalid format"),
		)
	})

	req := httptest.NewRequest("GET", "/error-with-fields", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 422 {
		t.Errorf("Status = %d, want 422", resp.StatusCode)
	}

	_ = l.Sync()
}

func TestErrorResponder_Error_EmptyLogKey(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})
	responder := NewErrorResponder(l)

	app := fiber.New()
	app.Use(requestid.New())

	app.Get("/error-no-log", func(c *fiber.Ctx) error {
		// Empty logKey should skip logging
		return responder.Error(c, 400, "Error without logging", "")
	})

	req := httptest.NewRequest("GET", "/error-no-log", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("Status = %d, want 400", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errorResp utils.ErrorResponseMap
	if err := json.Unmarshal(body, &errorResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if errorResp.Error != "Error without logging" {
		t.Errorf("Response.Error = %v, want 'Error without logging'", errorResp.Error)
	}

	_ = l.Sync()
}

func TestErrorResponder_Error_RequestIDPropagation(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})
	responder := NewErrorResponder(l)

	app := fiber.New()
	app.Use(requestid.New())

	app.Get("/error", func(c *fiber.Ctx) error {
		return responder.Error(c, 500, "Server error", "server_error")
	})

	customReqID := "test-request-id-456"
	req := httptest.NewRequest("GET", "/error", nil)
	req.Header.Set(fiber.HeaderXRequestID, customReqID)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var errorResp utils.ErrorResponseMap
	if err := json.Unmarshal(body, &errorResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if errorResp.RequestID != customReqID {
		t.Errorf("Response.RequestID = %v, want %v", errorResp.RequestID, customReqID)
	}

	_ = l.Sync()
}

func TestErrorResponder_Error_ContentType(t *testing.T) {
	l := logger.New(logger.Config{Debug: true})
	responder := NewErrorResponder(l)

	app := fiber.New()
	app.Use(requestid.New())

	app.Get("/error", func(c *fiber.Ctx) error {
		return responder.Error(c, 400, "Test error", "test")
	})

	req := httptest.NewRequest("GET", "/error", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %v, want application/json", contentType)
	}

	_ = l.Sync()
}
