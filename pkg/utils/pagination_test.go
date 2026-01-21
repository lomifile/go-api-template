package utils

import "testing"

func TestSortOptions(t *testing.T) {
	var opt SortOptions = "created_at"

	if string(opt) != "created_at" {
		t.Errorf("SortOptions = %v, want created_at", opt)
	}
}

func TestSortOptions_Assignment(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"created_at", "created_at"},
		{"updated_at", "updated_at"},
		{"name", "name"},
		{"id", "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opt SortOptions = SortOptions(tt.value)
			if string(opt) != tt.value {
				t.Errorf("SortOptions = %v, want %v", opt, tt.value)
			}
		})
	}
}
