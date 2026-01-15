package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// --------------------
// Unit Tests
// --------------------

func TestParseFile(t *testing.T) {

	tests := []struct {
		name string
		data string
		want map[string]string
	}{
		{
			name: "Standard Key-Value with Quotes",
			data: `key1 = "value1"
                   key2 = "value2"`,
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name: "Mixed Comments and spacing",
			data: `# This is a comment
                   // This is also a comment
                   
                   key_clean = "clean_value"
                   
                   spaced_key   =   "spaced_value"`,
			want: map[string]string{
				"key_clean":  "clean_value",
				"spaced_key": "spaced_value",
			},
		},
		{
			name: "Values without Quotes",
			data: `simple = no_quotes
                   number = 123`,
			want: map[string]string{
				"simple": "no_quotes",
				"number": "123",
			},
		},
		{
			name: "Invalid lines handling",
			data: `ThisLineHasNoEquals
                   =MissingKey
                   JustKey=`,
			want: map[string]string{"JustKey": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFile([]byte(tt.data))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "remora_test_*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `
host = "localhost"
port = "8080"
debug = "false"
`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	os.Setenv("port", "9090")
	os.Setenv("new_var", "should_be_ignored")
	defer os.Unsetenv("port")
	defer os.Unsetenv("new_var")

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	tests := []struct {
		key  string
		want string
	}{
		{"host", "localhost"}, // From file
		{"port", "9090"},      // Overridden by OS
		{"debug", "false"},    // From file
	}

	for _, tt := range tests {
		if got := cfg[tt.key]; got != tt.want {
			t.Errorf("Config key [%s] = %q, want %q", tt.key, got, tt.want)
		}
	}
	if _, exists := cfg["new_var"]; exists {
		t.Errorf("Expected 'new_var' to be ignored because it wasn't in the .env file, but it was found.")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {

	_, err := LoadConfig("non_existent_file_xyz.env")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestConfigReaders(t *testing.T) {

	cfg := ConfigMap{
		"str":      "hello",
		"bool_t":   "true",
		"bool_f":   "false",
		"bool_bad": "not_bool",
		"int":      "42",
		"int_bad":  "xyz",
		"flt":      "3.14",
		"flt_bad":  "abc",
	}

	t.Run("ReadString", func(t *testing.T) {
		if got := cfg.ReadString("str", "def"); got != "hello" {
			t.Errorf("want hello, got %s", got)
		}
		if got := cfg.ReadString("missing", "def"); got != "def" {
			t.Errorf("want def, got %s", got)
		}
	})

	t.Run("ReadBoolean", func(t *testing.T) {
		if got := cfg.ReadBoolean("bool_t", false); !got {
			t.Error("want true")
		}
		if _ = cfg.ReadBoolean("bool_bad", true); !true { // Fallback test
			t.Error("want fallback true for bad input")
		}
		if _ = cfg.ReadBoolean("missing", true); !true {
			t.Error("want fallback true for missing")
		}
	})

	t.Run("ReadInt", func(t *testing.T) {
		if got := cfg.ReadInt("int", 0); got != 42 {
			t.Errorf("want 42, got %d", got)
		}
		if got := cfg.ReadInt("int_bad", 10); got != 10 {
			t.Errorf("want fallback 10, got %d", got)
		}
	})

	t.Run("ReadFloat", func(t *testing.T) {
		if got := cfg.ReadFloat("flt", 0.0); got != 3.14 {
			t.Errorf("want 3.14, got %f", got)
		}
		if got := cfg.ReadFloat("flt_bad", 1.23); got != 1.23 {
			t.Errorf("want fallback 1.23, got %f", got)
		}
	})
}

// --------------------
// Benchmarks
// --------------------

func BenchmarkParseFile(b *testing.B) {

	data := []byte(`key1 = "value1"
key2 = "value2"
# comment
key3 = "value3"
key4 = "value4"`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parseFile(data)
	}
}

func BenchmarkParseFile_Large(b *testing.B) {

	var sb strings.Builder
	for i := range 1000 {
		sb.WriteString("key" + strconv.Itoa(i) + " = \"val" + strconv.Itoa(i) + "\"\n")
	}
	data := []byte(sb.String())

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parseFile(data)
	}
}

func BenchmarkLoadConfig(b *testing.B) {

	tmpFile, _ := os.CreateTemp("", "bench_*.env")
	defer os.Remove(tmpFile.Name())
	content := []byte(`a="1"\nb="2"\nc="3"`)
	_ = os.WriteFile(tmpFile.Name(), content, 0644)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := LoadConfig(tmpFile.Name())
		if err != nil {
			b.Fatalf("LoadConfig failed: %v", err)
		}
	}
}

func BenchmarkConfigReads(b *testing.B) {

	cfg := ConfigMap{
		"str":  "hello",
		"bool": "true",
		"int":  "42",
		"flt":  "3.14",
	}

	b.Run("ReadString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.ReadString("str", "fallback")
		}
	})
	b.Run("ReadBoolean", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.ReadBoolean("bool", false)
		}
	})
	b.Run("ReadInt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.ReadInt("int", 0)
		}
	})
}
