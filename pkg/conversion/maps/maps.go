package maps

import "strings"

// ExtractFieldValue returns the value associated with a field path in a map[string]any.
// The field path may contain dot notation to access nested maps.
//
// If any segment in the path does not exist, the function returns prematurely with a nil value.
//
// The returned value may be a primitive, a nested map, or any other type stored in the map.
func ExtractFieldValue(input map[string]any, field string) any {

	var value any = input
	fieldRemainder := field

	for fieldRemainder != "" {

		var fieldSection string

		pointIndex := strings.IndexByte(fieldRemainder, '.')
		if pointIndex >= 0 {

			fieldSection = fieldRemainder[:pointIndex]
			fieldRemainder = fieldRemainder[pointIndex+1:]

		} else {

			fieldSection = fieldRemainder
			fieldRemainder = ""
		}

		// If the input is not a map, it is the last section of the field
		castedValue, isMap := value.(map[string]any)
		if !isMap {
			return nil
		}

		// From partition of input map, try to get nested field
		nestedFieldValue, exists := castedValue[fieldSection]
		if !exists {
			return nil
		}
		value = nestedFieldValue
	}

	// Return the found nested field
	return value
}

// DeepCopy creates a fully independent deep copy of a map[string]any. All nested maps and slices
// are recursively copied so that modifications to the copy do not affect the original map.
// Primitive values and unhandled types (e.g., structs) are copied by reference.
//
// Returns nil if the original map is nil.
//
// Note: This function assumes that the map keys are strings and map values are of type any. It does not
// handle other types like structs or custom types. If such types are present as keys, they will be
// assigned directly without deep copying, may leading to shared references.
func DeepCopy(originalMap map[string]any) map[string]any {

	if originalMap == nil {
		return nil
	}

	copyMap := make(map[string]any)

	for key, value := range originalMap {

		switch valueType := value.(type) {

		// Nested map: perform a recursive deep copy
		case map[string]any:
			copyMap[key] = DeepCopy(valueType)

		// Slice: copy each element
		case []any:
			copySlice := make([]any, len(valueType))

			// For each item, if it is a map do a recursive deep copy
			for index, itemInSlice := range valueType {
				itemInSliceAsMap, isItemMap := itemInSlice.(map[string]any)
				if isItemMap {
					copySlice[index] = DeepCopy(itemInSliceAsMap)
				} else {
					copySlice[index] = itemInSlice
				}
			}
			copyMap[key] = copySlice

		// Primitive or unhandled type: assign directly
		default:
			copyMap[key] = value
		}
	}

	return copyMap
}

// ToLowerKeys returns a shallow copy of the given map[string]any with all keys converted to lowercase.
// The function works recursively: nested maps are converted with the same logic, and nested slices
// are processed recursively through sliceToLowerCase to ensure all keys in nested structures are normalized.
func ToLowerKeys(originalMap map[string]any) map[string]any {

	copiedMap := make(map[string]any)
	for key, value := range originalMap {

		lowerKey := strings.ToLower(key)

		switch nestedValue := value.(type) {

		case map[string]any:
			copiedMap[lowerKey] = ToLowerKeys(nestedValue)
		case []any:
			copiedMap[lowerKey] = sliceToLowerCase(nestedValue)
		default:
			copiedMap[lowerKey] = nestedValue
		}
	}
	return copiedMap
}

// sliceToLowerCase returns a shallow copy of the provided []any slice, recursively converting
// all keys in any nested maps to lowercase. Nested slices are processed recursively while
// primitive values are copied as-is.
func sliceToLowerCase(originalSlice []any) []any {

	copiedSlice := make([]any, len(originalSlice))
	for i, value := range originalSlice {

		switch nestedValue := value.(type) {

		case map[string]any:
			copiedSlice[i] = ToLowerKeys(nestedValue)
		case []any:
			copiedSlice[i] = sliceToLowerCase(nestedValue)
		default:
			copiedSlice[i] = nestedValue
		}
	}
	return copiedSlice
}
