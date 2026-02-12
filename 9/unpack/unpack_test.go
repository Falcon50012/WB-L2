package unpack

import (
	"errors"
	"strings"
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "basic unpack",
			input: "a4bc2d5e",
			want:  "aaaabccddddde",
		},
		{
			name:  "no digits",
			input: "abcd",
			want:  "abcd",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "zero repeat removes symbol",
			input: "a0",
			want:  "",
		},
		{
			name:  "multi-digit repeat",
			input: "a12",
			want:  strings.Repeat("a", 12),
		},
		{
			name:  "escaped digits",
			input: "qwe\\4\\5",
			want:  "qwe45",
		},
		{
			name:  "escaped digit then multiplier",
			input: "qwe\\45",
			want:  "qwe44444",
		},
		{
			name:  "escaped backslash",
			input: "\\\\",
			want:  "\\",
		},
		{
			name:  "complex example",
			input: "\\qwe\\45\\\\",
			want:  "qwe44444\\",
		},
		{
			name:    "starts with digit",
			input:   "45",
			wantErr: true,
		},
		{
			name:    "digit without preceding symbol",
			input:   "3abc",
			wantErr: true,
		},
		{
			name:    "dangling escape",
			input:   "\\",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got nil", tt.input)
				}
				if !errors.Is(err, ErrInvalidString) {
					t.Fatalf("expected ErrInvalidString, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}

			if got != tt.want {
				t.Fatalf("Unpack(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
