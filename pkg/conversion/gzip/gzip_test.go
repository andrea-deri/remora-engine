package gzip

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"strings"
	"testing"
)

// mustCompress is a test helper that performs the inverse operation of the function under test.
// It takes a plain string, compresses it using GZIP, and encodes the result in Base64.
//
// The function panics on any error, as failures are considered unrecoverable in test scenarios.
func mustCompress(input string) string {

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write([]byte(input)); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func TestCompressedStringToPlainString(t *testing.T) {

	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{name: "Valid Input - Simple String", input: mustCompress("Hello Remora"), expected: "Hello Remora", expectError: false},
		{name: "Valid Input - JSON Content", input: mustCompress(`{"status": "ok", "value": 123}`), expected: `{"status": "ok", "value": 123}`, expectError: false},
		{name: "Valid Input - Empty String", input: mustCompress(""), expected: "", expectError: false},
		{name: "Invalid Base64", input: "NotValidBase64!!", expected: "", expectError: true},
		{name: "Valid Base64 but Invalid Gzip", input: base64.StdEncoding.EncodeToString([]byte("This is not gzip data")), expected: "", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompressedStringToPlainString(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("CompressedStringToPlainString() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && got != tt.expected {
				t.Errorf("CompressedStringToPlainString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// --- Benchmarks ---

// Benchmark con payload piccolo (tipico di configurazioni o chiavi brevi)
func BenchmarkCompressedStringToPlainString_Small(b *testing.B) {
	input := mustCompress("small-config-string")

	b.ReportAllocs()

	for b.Loop() {
		_, _ = CompressedStringToPlainString(input)
	}
}

// Benchmark con payload medio (es. un JSON di risposta tipico, 2KB)
func BenchmarkCompressedStringToPlainString_Medium(b *testing.B) {

	payload := strings.Repeat("a", 2048)
	input := mustCompress(payload)

	b.ReportAllocs()

	for b.Loop() {
		_, _ = CompressedStringToPlainString(input)
	}
}
