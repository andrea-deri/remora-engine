package format

import (
	"testing"
)

// --------------------
// Unit Tests
// --------------------

func TestRemoveFirstChar(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Standard string",
			input: "/api/v1",
			want:  "api/v1",
		},
		{
			name:  "Two chars",
			input: "go",
			want:  "o",
		},
		{
			name:  "Single char",
			input: "a",
			want:  "",
		},
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Numeric string",
			input: "12345",
			want:  "2345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveFirstChar(tt.input); got != tt.want {
				t.Errorf("RemoveFirstChar(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestUnescape(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Double Quotes",
			input: `"value"`,
			want:  "value",
		},
		{
			name:  "Single Quotes",
			input: `'value'`,
			want:  "value",
		},
		{
			name:  "No Quotes",
			input: "value",
			want:  "value",
		},
		{
			name:  "Mismatched Quotes (Start)",
			input: `"value`,
			want:  `"value`,
		},
		{
			name:  "Mismatched Quotes (End)",
			input: `value"`,
			want:  `value"`,
		},
		{
			name:  "Mixed Quotes",
			input: `'value"`,
			want:  `'value"`,
		},
		{
			name:  "Nested Quotes (Double outer)",
			input: `"'inner'"`,
			want:  `'inner'`, // Trim removes outer ", keeps inner '
		},
		{
			name:  "Quotes inside value",
			input: `"val'ue"`,
			want:  `val'ue`,
		},
		{
			name:  "Empty String",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unescape(tt.input); got != tt.want {
				t.Errorf("Unescape(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// --------------------
// Benchmarks
// --------------------

func BenchmarkRemoveFirstChar(b *testing.B) {

	input := "/endpoints/user/1234"

	b.ReportAllocs()

	for b.Loop() {
		RemoveFirstChar(input)
	}
}

func BenchmarkUnescape(b *testing.B) {

	quoted := `"some_configuration_value"`
	unquoted := `simple_value`

	b.Run("Quoted", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Unescape(quoted)
		}
	})

	b.Run("Unquoted", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Unescape(unquoted)
		}
	})
}
