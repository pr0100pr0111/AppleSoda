package ui

import (
	"testing"
)

func TestTruncateStr(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hell…"},
		{"hello", 0, ""},
		{"hello", 1, "…"},
		{"ab", 2, "ab"},
		{"abc", 2, "a…"},
		{"", 5, ""},
	}
	for _, tt := range tests {
		got := truncateStr(tt.input, tt.max)
		if got != tt.expected {
			t.Errorf("truncateStr(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0:00"},
		{-1, "0:00"},
		{65, "1:05"},
		{3661, "61:01"},
		{30, "0:30"},
	}
	for _, tt := range tests {
		got := formatDuration(tt.input)
		if got != tt.expected {
			t.Errorf("formatDuration(%f) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
