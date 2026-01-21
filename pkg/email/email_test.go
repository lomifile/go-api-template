package email

import (
	"testing"
)

func TestSendEmailConfig(t *testing.T) {
	cfg := SendEmailConfig{
		To:                "test@example.com",
		Subject:           "Test Subject",
		AlternativeString: "Plain text content",
		HTML:              "<h1>HTML content</h1>",
	}

	if cfg.To != "test@example.com" {
		t.Errorf("To = %v, want test@example.com", cfg.To)
	}
	if cfg.Subject != "Test Subject" {
		t.Errorf("Subject = %v, want Test Subject", cfg.Subject)
	}
	if cfg.AlternativeString != "Plain text content" {
		t.Errorf("AlternativeString = %v, want Plain text content", cfg.AlternativeString)
	}
	if cfg.HTML != "<h1>HTML content</h1>" {
		t.Errorf("HTML = %v, want <h1>HTML content</h1>", cfg.HTML)
	}
}

func TestClient_Struct(t *testing.T) {
	// Test that Client struct can be created
	client := &Client{
		from: "sender@example.com",
	}

	if client.from != "sender@example.com" {
		t.Errorf("from = %v, want sender@example.com", client.from)
	}
}

func TestSendEmailConfig_EmptyFields(t *testing.T) {
	cfg := SendEmailConfig{}

	if cfg.To != "" {
		t.Errorf("To should be empty, got %v", cfg.To)
	}
	if cfg.Subject != "" {
		t.Errorf("Subject should be empty, got %v", cfg.Subject)
	}
	if cfg.AlternativeString != "" {
		t.Errorf("AlternativeString should be empty, got %v", cfg.AlternativeString)
	}
	if cfg.HTML != "" {
		t.Errorf("HTML should be empty, got %v", cfg.HTML)
	}
}
