package format

import "strings"

// RemoveFirstChar returns the input string without its first character.
// If the string is empty or has only one character, it returns an empty string.
func RemoveFirstChar(content string) string {

	if len(content) > 1 {
		return content[1:]
	}
	return ""
}

// Unescape removes surrounding single or double quotes from the provided value.
// If the value is not enclosed in matching quotes, it is returned unchanged.
//
// Example:
//
//	Input: "'someValue'"
//	Output: "someValue"
func Unescape(value string) string {

	if value == "" {
		return value
	}
	if strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`) {
		return value[1 : len(value)-1]
	}
	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		return value[1 : len(value)-1]
	}
	return value
}
