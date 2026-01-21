package utils

import "testing"

func TestSQLOrderTypes(t *testing.T) {
	tests := []struct {
		name     string
		order    SQLOrderTypes
		expected string
	}{
		{"Ascending", Asc, "ASC"},
		{"Descending", Desc, "DESC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.order) != tt.expected {
				t.Errorf("SQLOrderTypes = %v, want %v", tt.order, tt.expected)
			}
		})
	}
}

func TestSQLOrderTypes_String(t *testing.T) {
	asc := Asc
	desc := Desc

	if string(asc) != "ASC" {
		t.Errorf("Asc = %v, want ASC", asc)
	}
	if string(desc) != "DESC" {
		t.Errorf("Desc = %v, want DESC", desc)
	}
}
