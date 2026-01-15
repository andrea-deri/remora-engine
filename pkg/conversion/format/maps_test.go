package format

import (
	"reflect"
	"testing"
)

// --------------------
// Unit Tests
// --------------------

func TestFromAnyToStringMap(t *testing.T) {

	tests := []struct {
		name  string
		input map[string]any
		want  map[string]string
	}{
		{
			name:  "Nil Input",
			input: nil,
			want:  map[string]string{},
		},
		{
			name:  "Empty Input",
			input: map[string]any{},
			want:  map[string]string{},
		},
		{
			name: "Basic Types (String, Int, Bool)",
			input: map[string]any{
				"key_str":  "hello",
				"key_int":  123,
				"key_bool": true,
			},
			want: map[string]string{
				"key_str":  "hello",
				"key_int":  "123",
				"key_bool": "true",
			},
		},
		{
			name: "Float Types",
			input: map[string]any{
				"pi":    3.14,
				"whole": 10.0,
			},
			want: map[string]string{
				"pi":    "3.14",
				"whole": "10",
			},
		},
		{
			name: "Complex Types (Slices, Maps)",
			input: map[string]any{
				"slice": []int{1, 2},
				"map":   map[string]int{"a": 1},
			},
			want: map[string]string{
				"slice": "[1 2]",
				"map":   "map[a:1]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromAnyToStringMap(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromAnyToStringMap() = %v, want %v", got, tt.want)
			}
			if got == nil {
				t.Error("FromAnyToStringMap() returned nil map")
			}
		})
	}
}

// --------------------
// Benchmarks
// --------------------

func BenchmarkFromAnyToStringMap_Small(b *testing.B) {

	input := map[string]any{
		"Content-Type": "application/json",
		"Retries":      3,
		"Is-Active":    true,
		"Latency":      0.45,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		FromAnyToStringMap(input)
	}
}

func BenchmarkFromAnyToStringMap_Large(b *testing.B) {

	input := make(map[string]any, 100)
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			input["k"+string(rune(i))] = i
		} else {
			input["k"+string(rune(i))] = "value"
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		FromAnyToStringMap(input)
	}
}
