package maps

import (
	"reflect"
	"testing"
)

// --------------------
// Unit Tests
// --------------------

func TestExtractFieldValue(t *testing.T) {

	data := map[string]any{
		"simple": "value",
		"level1": map[string]any{
			"level2": map[string]any{
				"target": "bingo",
				"value":  42,
			},
			"leaf": "end",
		},
		"integers": 100,
	}

	tests := []struct {
		name  string
		input map[string]any
		field string
		want  any
	}{
		{name: "Top level string", input: data, field: "simple", want: "value"},
		{name: "Nested deep access", input: data, field: "level1.level2.target", want: "bingo"},
		{name: "Nested deep integer", input: data, field: "level1.level2.value", want: 42},
		{name: "Missing Key Top Level", input: data, field: "missing", want: nil},
		{name: "Missing Key Nested", input: data, field: "level1.missing", want: nil},
		{name: "Path exists but is not a map (Stop early)", input: data, field: "integers.subfield", want: nil},
		{name: "Nil input map", input: nil, field: "anything", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractFieldValue(tt.input, tt.field)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractFieldValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeepCopy(t *testing.T) {

	t.Run("Should return nil for nil input", func(t *testing.T) {
		if DeepCopy(nil) != nil {
			t.Error("Expected nil return for nil input")
		}
	})

	t.Run("Should copy recursively and independently", func(t *testing.T) {

		original := map[string]any{
			"scalar": 1,
			"nested": map[string]any{
				"inner": "original",
			},
			"slice": []any{
				map[string]any{"id": 1},
				"simple_string",
			},
		}

		copyMap := DeepCopy(original)

		if !reflect.DeepEqual(original, copyMap) {
			t.Fatal("DeepCopy content differs from original")
		}

		nestedMap := copyMap["nested"].(map[string]any)
		nestedMap["inner"] = "modified"

		if original["nested"].(map[string]any)["inner"] == "modified" {
			t.Error("Modifying copy affected the original map (Nested Map)")
		}

		slice := copyMap["slice"].([]any)
		sliceMap := slice[0].(map[string]any)
		sliceMap["id"] = 999

		originalSlice := original["slice"].([]any)
		originalSliceMap := originalSlice[0].(map[string]any)

		if originalSliceMap["id"] == 999 {
			t.Error("Modifying copy affected the original map (Slice Content)")
		}
	})
}

func TestToLowerKeys(t *testing.T) {

	input := map[string]any{
		"KeyOne": "Value1",
		"NESTED": map[string]any{
			"InnerKey": "InnerValue",
		},
		"List": []any{
			map[string]any{"ListKey": 1},
			"StringItem",
		},
	}

	want := map[string]any{
		"keyone": "Value1",
		"nested": map[string]any{
			"innerkey": "InnerValue",
		},
		"list": []any{
			map[string]any{"listkey": 1},
			"StringItem",
		},
	}

	got := ToLowerKeys(input)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToLowerKeys() mismatch.\nGot:  %v\nWant: %v", got, want)
	}
	if _, exists := input["KeyOne"]; !exists {
		t.Error("ToLowerKeys modified the original map in place, expected shallow copy behavior for top level")
	}
}

// --------------------
// Benchmarks
// --------------------

func BenchmarkExtractFieldValue(b *testing.B) {

	data := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{
					"d": "found",
				},
			},
		},
	}

	b.Run("Shallow", func(b *testing.B) {
		for b.Loop() {
			ExtractFieldValue(data, "a")
		}
	})

	b.Run("Deep", func(b *testing.B) {
		for b.Loop() {
			ExtractFieldValue(data, "a.b.c.d")
		}
	})

	b.Run("Deep_ReportAllocs", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			ExtractFieldValue(data, "a.b.c.d")
		}
	})
}

func BenchmarkDeepCopy(b *testing.B) {

	smallMap := map[string]any{"a": 1, "b": "2"}

	largeMap := make(map[string]any)
	for i := range 100 {
		largeMap["key"] = map[string]any{
			"nested": i,
			"data":   make([]any, 10),
		}
	}

	b.Run("Small", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			DeepCopy(smallMap)
		}
	})

	b.Run("Large", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			DeepCopy(largeMap)
		}
	})
}

func BenchmarkToLowerKeys(b *testing.B) {

	input := map[string]any{
		"CamelCase": "val",
		"NESTED": map[string]any{
			"DEEP": "val",
		},
	}

	b.ReportAllocs()
	for b.Loop() {
		ToLowerKeys(input)
	}
}
