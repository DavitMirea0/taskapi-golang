package service

import (
	"errors"
	"strings"
	"testing"

	"taskapi/internal/domain"
)

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		want    string
		wantErr bool
	}{
		{"ok", "buy milk", "buy milk", false},
		{"trims spaces", "  buy milk  ", "buy milk", false},
		{"empty", "", "", true},
		{"only spaces", "   ", "", true},
		{"too long", strings.Repeat("a", 201), "", true},
		{"max length", strings.Repeat("a", 200), strings.Repeat("a", 200), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTitle(tt.title)
			if tt.wantErr {
				if !errors.Is(err, domain.ErrValidation) {
					t.Fatalf("want ErrValidation, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
