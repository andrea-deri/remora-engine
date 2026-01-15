package resolution

import (
	"fmt"
	"regexp"
	"remora/pkg/conversion"
	"strings"
	"unicode"
)

// QueueField represents a node in the breadth-first traversal queue used when
// scanning nested map and slice structures.
type QueueField struct {

	// value is the current element being inspected
	value any

	// fieldPathPrefix stores the full dot-notation path within the root object.
	fieldPathPrefix string
}

const (
	// RESOLUTION_DIRECTIVE_CASTTOBOOLEAN identifies a resolution directive
	// that converts a placeholder’s string value into a boolean type.
	RESOLUTION_DIRECTIVE_CASTTOBOOLEAN = "cbo"

	// RESOLUTION_DIRECTIVE_CASTTONUMBER identifies a resolution directive
	// that converts a placeholder’s string value into a numeric type.
	RESOLUTION_DIRECTIVE_CASTTONUMBER = "cnu"
)

// Regular expression to match placeholders in the format "${...}"
var resolvableValueRegex = regexp.MustCompile(`\$\{(.*?)\}`)

// InsertResolvedValueInMappableBody sets a value inside a nested map[string]any structure using a
// dot-notated path.
// The function guarantees that the full path exists by creating intermediate maps when missing.
// If a non-map value is encountered along the path, it is replaced with a new map to allow
// continuation of the descent.
//
// On the final field:
//   - If the existing value is a string, placeholder substitution is performed using the pattern
//     "${parameter}" and resolution directives (e.g. type casting) are applied when present.
//   - For all other types, the value is directly assigned.
//
// The function mutates the input map in-place.
func InsertResolvedValueInMappableBody(mappableBody map[string]any, stringToResolve string, fieldPath string, valueToSet any) {

	pathSections := strings.Split(fieldPath, ".")
	lastPathSectionIndex := len(pathSections) - 1
	mappableBodySubmap := mappableBody

	for i, fieldFromPathSection := range pathSections {

		if i == lastPathSectionIndex {

			switch valueFromField := mappableBodySubmap[fieldFromPathSection].(type) {

			case string:
				resolvedValueByReplace := strings.ReplaceAll(valueFromField, fmt.Sprintf("${%s}", stringToResolve), fmt.Sprintf("%v", valueToSet))
				resolvedValueByDirective := applyResolutionDirective(valueFromField, resolvedValueByReplace)
				mappableBodySubmap[fieldFromPathSection] = resolvedValueByDirective

			default:
				mappableBodySubmap[fieldFromPathSection] = valueToSet
			}

			return
		}

		valueOfFieldFromPathSection, existsFieldAtPath := mappableBodySubmap[fieldFromPathSection]
		if !existsFieldAtPath {

			// No field found at path: add it as new node
			newValueOfFieldFromPathSection := make(map[string]any)
			mappableBodySubmap[fieldFromPathSection] = newValueOfFieldFromPathSection
			mappableBodySubmap = newValueOfFieldFromPathSection
			continue
		}

		inputSubMap, isValueMap := valueOfFieldFromPathSection.(map[string]any)
		if !isValueMap {

			// No map field found at path: add it as new node
			newValueOfFieldFromPathSection := make(map[string]any)
			mappableBodySubmap[fieldFromPathSection] = newValueOfFieldFromPathSection
			mappableBodySubmap = newValueOfFieldFromPathSection
			continue
		}

		// Jump to next hierarchy level
		mappableBodySubmap = inputSubMap
	}
}

// GetResolvableFieldsFromString scans the input string for placeholders matching the ${...} pattern.
// It returns a slice containing all extracted placeholder names (the inner content without ${}).
// If no placeholders are found, it returns an empty slice.
func GetResolvableFieldsFromString(content string) []string {

	resolvablePlaceholdersOnContent := resolvableValueRegex.FindAllStringSubmatch(content, -1)
	if resolvablePlaceholdersOnContent == nil {
		return []string{}
	}

	resolvableFieldNames := make([]string, len(resolvablePlaceholdersOnContent))
	for i, resolvablePlaceholder := range resolvablePlaceholdersOnContent {
		resolvableFieldNames[i] = resolvablePlaceholder[1]
	}

	return resolvableFieldNames
}

// GetResolvableFieldNamesFromMap scans a nested map[string]any structure and returns the list of all
// field paths (in dot notation) whose string values contain placeholders.
// A placeholder is any substring matching the pattern defined by resolvableValueRegex.
//
// The traversal is breadth-first and supports:
//   - nested maps (paths are expanded using dot notation)
//   - slices (each element is inspected, preserving the parent path)
//   - strings (only those containing a placeholder produce an output path)
//
// The function returns only the paths of fields whose final value is a string
// that includes at least one placeholder.
//
// Example:
//
//	Input:  data["path"]["to"]["field"] = "${body.param}"
//	Output: ["path.to.field"]
func GetResolvableFieldNamesFromMap(data map[string]any) []string {

	// Slice collecting the dot-notation paths of all fields whose string value
	// contains a placeholder. Paths are discovered via breadth-first traversal.
	var extractedParameters []string

	// Initialize the queue: first node to analyze is the entire map
	fieldsToAnalyzeQueue := []QueueField{{value: data, fieldPathPrefix: ""}}

	for len(fieldsToAnalyzeQueue) > 0 {

		// Pop item to analyze from queue
		itemToAnalyze := fieldsToAnalyzeQueue[0]
		fieldsToAnalyzeQueue = fieldsToAnalyzeQueue[1:]

		switch castedItemToAnalyze := itemToAnalyze.value.(type) {

		case map[string]any:
			// If the value is a map, enqueue each key-value pair with complete path
			for itemKey, itemValue := range castedItemToAnalyze {
				completePath := itemKey
				if itemToAnalyze.fieldPathPrefix != "" {
					completePath = itemToAnalyze.fieldPathPrefix + "." + completePath
				}
				fieldsToAnalyzeQueue = append(fieldsToAnalyzeQueue, QueueField{value: itemValue, fieldPathPrefix: completePath})
			}

		case []any:
			// If the value is a slice, enqueue each item without changing the prefix
			for _, itemFromSlice := range castedItemToAnalyze {
				fieldsToAnalyzeQueue = append(fieldsToAnalyzeQueue, QueueField{value: itemFromSlice, fieldPathPrefix: itemToAnalyze.fieldPathPrefix})
			}

		case string:
			if resolvableValueRegex.MatchString(castedItemToAnalyze) {
				extractedParameters = append(extractedParameters, itemToAnalyze.fieldPathPrefix)
			}
		}
	}

	return extractedParameters
}

// GetResolvableFieldNamesFromFunction parses a function-like expression and extracts both the function
// name and its argument list. The expected input format is:
//
//	!functionName(param1, param2, ...)
//
// The parser supports:
//   - Zero or more parameters
//   - Arbitrary whitespace
//   - Nested function calls (e.g. !op(!inner(a,b), 1))
//   - Arbitrary expressions as parameters, including literals and placeholders as '$value.path',
//   - Proper handling of balanced parentheses to avoid incorrect splits on commas inside nested structures
//
// The returned values are:
//   - the extracted function name (string)
//   - a slice of parameters in their raw textual form, without modification
//     (e.g. ["a", "$body.val1", "!op2(!op3(),$qparam.val2)"])
//
// Behavior and edge cases:
//   - If the input does not contain a leading '!', the function returns ("", []string{}).
//   - If the string contains a function name but no parentheses, it returns the name and an empty parameter list.
//   - If parentheses are present but empty (e.g. "!op()"), it returns the name and an empty slice.
//   - Unbalanced parentheses inside the parameter list stop parsing at the first valid closure of the main function.
//
// Example:
//
//	Input:  "!randomInt(1, !div(10, 2))"
//	Output: functionName = "randomInt"
//	        parameters   = ["1", "!div(10, 2)"]
//
// This function is designed for deterministic extraction only: it does not evaluate or interpret
// any parameter content.
func GetResolvableFieldNamesFromFunction(resolvableFunctionPlaceholder string) (string, []string) {

	// Check function validity
	resolvableFunctionStartIndex := strings.Index(resolvableFunctionPlaceholder, "!")
	if resolvableFunctionStartIndex == -1 {
		return "", []string{}
	}
	resolutionFunction := resolvableFunctionPlaceholder[resolvableFunctionStartIndex+1:]

	functionNameEndIndex := strings.Index(resolutionFunction, "(")
	if functionNameEndIndex == -1 {
		return resolutionFunction, []string{}
	}
	functionName := resolutionFunction[:functionNameEndIndex]

	// Get the parameters after the first "("
	functionParametersRawString := resolutionFunction[functionNameEndIndex+1:]
	if functionParametersRawString == ")" {
		return functionName, []string{}
	}

	var parameters []string
	var newParameterBuffer strings.Builder
	var currentNestedFunctionCursor int
	var foundRawStringEnd bool

	// Executing raw parameter string parsing by parenthesis balancing
	for _, rawStringRune := range functionParametersRawString {

		if foundRawStringEnd {
			break
		}

		switch rawStringRune {

		case '(':
			// Open parenthesis for opening nested function
			currentNestedFunctionCursor++
			newParameterBuffer.WriteRune(rawStringRune)

		case ')':
			if currentNestedFunctionCursor > 0 {
				// Closure parenthesis for closing nested function
				currentNestedFunctionCursor--
				newParameterBuffer.WriteRune(rawStringRune)
			} else {
				// Closure parenthesis for closing main function: end computation
				foundRawStringEnd = true
			}

		case ',':
			if currentNestedFunctionCursor == 0 {
				// Main function: generate parameterFromBuffer from buffer
				parameterFromBuffer := strings.TrimSpace(newParameterBuffer.String())
				if parameterFromBuffer != "" {
					parameters = append(parameters, parameterFromBuffer)
				}
				newParameterBuffer.Reset()
			} else {
				// Nested function: write on buffer without generating sub-parameter
				newParameterBuffer.WriteRune(rawStringRune)
			}

		// For all other characters, add them on buffer
		default:
			if !unicode.IsControl(rawStringRune) {
				newParameterBuffer.WriteRune(rawStringRune)
			}
		}
	}

	// Generate last parameter if there is remaining content in buffer
	if newParameterBuffer.Len() > 0 {
		parameterFromBuffer := strings.TrimSpace(newParameterBuffer.String())
		if parameterFromBuffer != "" {
			parameters = append(parameters, parameterFromBuffer)
		}
	}

	return functionName, parameters
}

// ExtractFieldsFromResolutionPlaceholders scans the input string for parameter placeholders in the form "${...}"
// and extracts the internal field paths. This is typically used to identify which dynamic
// values must be resolved at runtime when processing templates, resolution expressions
// or rule-based configurations.
//
// The function returns only the raw parameter identifiers without delimiters, preserving
// their hierarchical dot-notation (e.g., "a.b.c"). If the input contains no placeholders,
// an empty slice is returned.
//
// Example:
//
//	Input:  "${path.to.param} = ${other.path}"
//	Output: ["path.to.param", "other.path"]
func ExtractFieldsFromResolutionPlaceholders(resolvableField string) []string {

	resolutionPlaceholders := resolvableValueRegex.FindAllStringSubmatch(resolvableField, -1)

	var resolvableFields []string
	for _, resolutionPlaceholder := range resolutionPlaceholders {
		if len(resolutionPlaceholder) > 1 {
			resolvableFields = append(resolvableFields, resolutionPlaceholder[1])
		}
	}
	return resolvableFields
}

// applyResolutionDirective processes a resolution directive embedded in the input string.
// A directive is identified by the "xxx::" prefix, where "xxx" specifies an operation
// to apply to the resolvable value.
//
// Supported directives:
//   - cnu: cast to number
//   - cbo: cast to boolean
//
// If the input string does not contain a valid directive prefix, the function returns
// the resolvableValue without modification.
func applyResolutionDirective(checkedValue string, resolvableValue any) any {

	// Return prematurely if no resolution directive (content with prefix 'xxx::') is found
	if len(checkedValue) <= 5 || checkedValue[3:5] != "::" {
		return resolvableValue
	}

	switch checkedValue[:3] {

	case RESOLUTION_DIRECTIVE_CASTTONUMBER:
		castedResolvedValue := resolvableValue.(string)[5:]
		resolvedValueAsInteger, isConversionMade := conversion.AsNumber(castedResolvedValue)
		if isConversionMade {
			return resolvedValueAsInteger
		}

	case RESOLUTION_DIRECTIVE_CASTTOBOOLEAN:
		castedResolvedValue := resolvableValue.(string)[5:]
		resolvedValueAsBoolean, isConversionMade := conversion.AsBoolean(castedResolvedValue)
		if isConversionMade {
			return resolvedValueAsBoolean
		}
	}

	return resolvableValue
}
