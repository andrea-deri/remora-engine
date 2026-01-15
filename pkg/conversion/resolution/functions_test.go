package resolution

import (
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
)

// --------------------
// Unit Tests
// --------------------

func TestFromFunction_Deterministic(t *testing.T) {

	tests := []struct {
		name     string
		funcName string
		params   []any
		want     any
	}{
		{name: "Add Integers", funcName: RESOLUTION_FUNCTION_ADD, params: []any{10.0, 5.0}, want: 15.0},
		{name: "Sub Floats", funcName: RESOLUTION_FUNCTION_SUB, params: []any{10.5, 0.5}, want: 10.0},
		{name: "Mul Mixed", funcName: RESOLUTION_FUNCTION_MUL, params: []any{2.0, 4.0}, want: 8.0},
		{name: "Div Valid", funcName: RESOLUTION_FUNCTION_DIV, params: []any{20.0, 4.0}, want: 5.0},
		{name: "Div By Zero (Should handle gracefully)", funcName: RESOLUTION_FUNCTION_DIV, params: []any{10.0, 0.0}, want: 0},
		{name: "Concat Strings", funcName: RESOLUTION_FUNCTION_CONCAT, params: []any{"Hello", "World"}, want: "HelloWorld"},
		{name: "Concat Invalid Types", funcName: RESOLUTION_FUNCTION_CONCAT, params: []any{123, "World"}, want: ""},
		{name: "Format Float Standard", funcName: RESOLUTION_FUNCTION_FORMATFLOAT, params: []any{3.14159, 2.0}, want: "3.14"},
		{name: "Format Float No Decimals", funcName: RESOLUTION_FUNCTION_FORMATFLOAT, params: []any{10.555, 0.0}, want: "11"},
		{name: "Unknown Function", funcName: "non_existent_function", params: nil, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromFunction(tt.funcName, tt.params)

			if gotNum, ok := got.(float64); ok {
				if wantNum, ok := tt.want.(float64); ok {
					if gotNum != wantNum {
						t.Errorf("FromFunction() = %v, want %v", got, tt.want)
					}
					return
				} else if wantNum, ok := tt.want.(int); ok {
					if int(gotNum) != wantNum {
						t.Errorf("FromFunction() = %v, want %v", got, tt.want)
					}
					return
				}
			}

			if got != tt.want {
				t.Errorf("FromFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromFunction_RandomAndGenerative(t *testing.T) {

	t.Run("UUID Generation", func(t *testing.T) {
		val := FromFunction(RESOLUTION_FUNCTION_UUID, nil)
		strVal, ok := val.(string)
		if !ok {
			t.Fatalf("Expected string for UUID, got %T", val)
		}
		if _, err := uuid.Parse(strVal); err != nil {
			t.Errorf("Generated UUID is invalid: %s", strVal)
		}
	})

	t.Run("Random Int Range", func(t *testing.T) {
		min, max := 10.0, 20.0
		val := FromFunction(RESOLUTION_FUNCTION_RANDOMINT, []any{min, max})
		intVal, ok := val.(int64)
		if !ok {
			t.Fatalf("Expected int64, got %T", val)
		}
		if intVal < int64(min) || intVal > int64(max) {
			t.Errorf("RandomInt %d out of bounds [%v, %v]", intVal, min, max)
		}
	})

	t.Run("Random String Length", func(t *testing.T) {
		length := 16.0
		val := FromFunction(RESOLUTION_FUNCTION_RANDOMSTRING, []any{length})
		strVal, ok := val.(string)
		if !ok {
			t.Fatalf("Expected string, got %T", val)
		}
		if len(strVal) != int(length) {
			t.Errorf("RandomString length = %d, want %d", len(strVal), int(length))
		}
	})
}

func TestFromFunction_Time(t *testing.T) {

	t.Run("Timestamp", func(t *testing.T) {
		now := time.Now().Unix()
		val := FromFunction(RESOLUTION_FUNCTION_TIMESTAMP, nil)
		tsVal, ok := val.(int64)
		if !ok {
			t.Fatalf("Expected int64, got %T", val)
		}
		// Tollerance of 2 seconds
		if math.Abs(float64(tsVal-now)) > 2 {
			t.Errorf("Timestamp deviation too high: got %d, expected approx %d", tsVal, now)
		}
	})

	t.Run("Today Format", func(t *testing.T) {
		val := FromFunction(RESOLUTION_FUNCTION_TODAY, nil)
		strVal, ok := val.(string)
		if !ok {
			t.Fatal("Expected string")
		}
		matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, strVal)
		if !matched {
			t.Errorf("Today format mismatch: %s", strVal)
		}
	})

	t.Run("Now Format", func(t *testing.T) {
		val := FromFunction(RESOLUTION_FUNCTION_NOW, nil)
		strVal, ok := val.(string)
		if !ok {
			t.Fatal("Expected string")
		}
		matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, strVal)
		if !matched {
			t.Errorf("Now format mismatch: %s. WARNING: Check Go time format string in implementation!", strVal)
		}
	})
}

// --------------------
// Benchmarks
// --------------------

func BenchmarkResolution_Add(b *testing.B) {

	params := []any{100.0, 200.0}
	b.ReportAllocs()
	for b.Loop() {
		FromFunction(RESOLUTION_FUNCTION_ADD, params)
	}
}

func BenchmarkResolution_UUID(b *testing.B) {

	b.ReportAllocs()
	for b.Loop() {
		FromFunction(RESOLUTION_FUNCTION_UUID, nil)
	}
}

// Benchmark per Random String (allocazione memoria)
func BenchmarkResolution_RandomString(b *testing.B) {

	params := []any{32.0}

	b.ReportAllocs()
	for b.Loop() {
		FromFunction(RESOLUTION_FUNCTION_RANDOMSTRING, params)
	}
}

// Benchmark per Format Float (conversione stringa)
func BenchmarkResolution_FormatFloat(b *testing.B) {

	params := []any{123.456789, 2.0}

	b.ReportAllocs()
	for b.Loop() {
		FromFunction(RESOLUTION_FUNCTION_FORMATFLOAT, params)
	}
}
