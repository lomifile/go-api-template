package utils

import (
	"encoding/hex"
	"testing"
)

func TestRandToken(t *testing.T) {
	tests := []struct {
		name       string
		byteLength int
		wantHexLen int
	}{
		{"8 bytes", 8, 16},
		{"16 bytes", 16, 32},
		{"32 bytes", 32, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := RandToken(tt.byteLength)
			if err != nil {
				t.Errorf("RandToken(%d) returned error: %v", tt.byteLength, err)
			}
			if len(token) != tt.wantHexLen {
				t.Errorf("RandToken(%d) returned token of length %d, want %d", tt.byteLength, len(token), tt.wantHexLen)
			}
		})
	}
}

func TestRandToken_ValidHex(t *testing.T) {
	token, err := RandToken(16)
	if err != nil {
		t.Fatalf("RandToken failed: %v", err)
	}

	_, err = hex.DecodeString(token)
	if err != nil {
		t.Errorf("RandToken returned invalid hex string: %v", err)
	}
}

func TestRandToken_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		token, err := RandToken(16)
		if err != nil {
			t.Fatalf("RandToken failed: %v", err)
		}
		if seen[token] {
			t.Errorf("RandToken produced duplicate value: %s", token)
		}
		seen[token] = true
	}
}
