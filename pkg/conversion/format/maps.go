package format

import (
	"fmt"
	"strconv"
)

// FromAnyToStringMap converts a map[string]any into a map[string]string by
// converting each value to its string representation.
//
// Primitive types are handled explicitly for efficiency, while unsupported
// types fall back to fmt.Sprintf.
func FromAnyToStringMap(input map[string]any) map[string]string {

	output := make(map[string]string, len(input))
	for key, value := range input {
		switch v := value.(type) {
		case string:
			output[key] = v
		case int:
			output[key] = strconv.Itoa(v)
		case bool:
			output[key] = strconv.FormatBool(v)
		case float64:
			output[key] = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			// Slower fallback only for different types
			output[key] = fmt.Sprintf("%v", v)
		}
	}
	return output
}
