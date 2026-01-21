package utils

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 1", 1},
		{"length 10", 10},
		{"length 32", 32},
		{"length 64", 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomString(tt.length)
			if len(result) != tt.length {
				t.Errorf("RandomString(%d) returned string of length %d, want %d", tt.length, len(result), tt.length)
			}
		})
	}
}

func TestRandomString_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	iterations := 1000
	length := 16

	for i := 0; i < iterations; i++ {
		s := RandomString(length)
		if seen[s] {
			t.Errorf("RandomString produced duplicate value: %s", s)
		}
		seen[s] = true
	}
}

func TestRandomString_Characters(t *testing.T) {
	// base64 URL encoding uses A-Z, a-z, 0-9, -, _
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	validSet := make(map[rune]bool)
	for _, c := range validChars {
		validSet[c] = true
	}

	s := RandomString(100)
	for _, c := range s {
		if !validSet[c] {
			t.Errorf("RandomString contains invalid character: %c", c)
		}
	}
}
