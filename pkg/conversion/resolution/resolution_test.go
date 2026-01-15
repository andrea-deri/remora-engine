package resolution

import (
	"reflect"
	"testing"
)

func TestGetParametersFromFunction(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantName   string
		wantParams []string
	}{
		{"NoParameters", "!op()", "op", []string{}},
		{"TwoSimpleParameters", "!op(a,b)", "op", []string{"a", "b"}},
		{"SimpleParameterAndPlaceholderParameter", "!op($body.val1,1)", "op", []string{"$body.val1", "1"}},
		{"TwoSimpleParametersAndPlaceholderParameter", "!op(3,$body.val1,1)", "op", []string{"3", "$body.val1", "1"}},
		{"NestedFunction", "!op(!op2(a,b),1)", "op", []string{"!op2(a,b)", "1"}},
		{"NestedFunctionWithPlaceholderParameter", "!op(1,!op2(a,$body.val1),1)", "op", []string{"1", "!op2(a,$body.val1)", "1"}},
		{"ThreeLayerNestedFunctionWithPlaceholderParameter", "!op(1,!op2(!op3(),$body.val1),!op4(ex))", "op", []string{"1", "!op2(!op3(),$body.val1)", "!op4(ex)"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotParams := GetResolvableFieldNamesFromFunction(tt.input)
			if gotName != tt.wantName {
				t.Errorf("GetParametersFromFunction() name = %v, want %v", gotName, tt.wantName)
			}
			if !reflect.DeepEqual(gotParams, tt.wantParams) {
				t.Errorf("GetParametersFromFunction() params = %v, want %v", gotParams, tt.wantParams)
			}
		})
	}
}
