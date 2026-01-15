package resolution

import (
	"math"
	"math/rand/v2"
	"strconv"
	"time"
	"unsafe"

	"remora/pkg/conversion"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (

	// RESOLUTION_FUNCTION_UUID represents the name of the function that generates a UUID string.
	RESOLUTION_FUNCTION_UUID string = "uuid"

	// RESOLUTION_FUNCTION_RANDOMINT represents the name of the function that generates a random integer.
	RESOLUTION_FUNCTION_RANDOMINT string = "randomInt"

	// RESOLUTION_FUNCTION_RANDOMFLOAT represents the name of the function that generates a random floating-point number.
	RESOLUTION_FUNCTION_RANDOMFLOAT string = "randomFloat"

	// RESOLUTION_FUNCTION_RANDOMSTRING represents the name of the function that generates a random alphanumeric string.
	RESOLUTION_FUNCTION_RANDOMSTRING string = "randomString"

	// RESOLUTION_FUNCTION_CONCAT represents the name of the function that concatenates multiple values into a single string.
	RESOLUTION_FUNCTION_CONCAT string = "concat"

	// RESOLUTION_FUNCTION_ADD represents the name of the function that performs numeric addition.
	RESOLUTION_FUNCTION_ADD string = "add"

	// RESOLUTION_FUNCTION_SUB represents the name of the function that performs numeric subtraction.
	RESOLUTION_FUNCTION_SUB string = "sub"

	// RESOLUTION_FUNCTION_MUL represents the name of the function that performs numeric multiplication.
	RESOLUTION_FUNCTION_MUL string = "mul"

	// RESOLUTION_FUNCTION_DIV represents the name of the function that performs numeric division.
	RESOLUTION_FUNCTION_DIV string = "div"

	// RESOLUTION_FUNCTION_FORMATFLOAT represents the name of the function that formats a float using a specific pattern.
	RESOLUTION_FUNCTION_FORMATFLOAT string = "formatFloat"

	// RESOLUTION_FUNCTION_NOW represents the name of the function that returns the current date and time.
	RESOLUTION_FUNCTION_NOW string = "now"

	// RESOLUTION_FUNCTION_TODAY represents the name of the function that returns the current date without time.
	RESOLUTION_FUNCTION_TODAY string = "today"

	// RESOLUTION_FUNCTION_TIMESTAMP represents the name of the function that returns the current Unix timestamp.
	RESOLUTION_FUNCTION_TIMESTAMP string = "timestamp"
)

// ALPHANUMERIC_BASIC_CHARSET defines the basic alphanumeric character set used for random string generation.
// It includes lowercase letters, uppercase letters, and digits.
const ALPHANUMERIC_BASIC_CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// FromFunction executes the specified function with the given parameters and returns the resolved value.
// Supported functions include:
//   - uuid(): Generates a new UUID string.
//   - randomInt(min, max): Generates a random integer between min and max (inclusive).
//   - randomFloat(min, max): Generates a random float between min and max.
//   - randomString(length): Generates a random alphanumeric string of the specified length.
//   - concat(prefix, suffix): Generates the concatenation of two strings.
//   - add(addend1, addend2): Returns the sum of two numbers.
//   - sub(minuend, subtrahend): Returns the difference between two numbers.
//   - mul(factor1, factor2): Returns the product of two numbers.
//   - div(numerator, denominator): Returns the quotient of two numbers (denominator must not be zero).
//   - formatFloat(number, precision): Formats a float to the specified number of decimal places.
//   - now(): Returns the current date and time as a formatted string "YYYY-MM-DD HH:MM:SS".
//   - today(): Returns the current date as a formatted string "YYYY-MM-DD".
//   - timestamp(): Returns the current Unix timestamp as an integer.
//
// If the function name is not recognized or parameters are invalid, it logs an error and returns nil.
func FromFunction(functionName string, parameters []any) any {

	switch functionName {

	case RESOLUTION_FUNCTION_UUID:
		return fromUUID()
	case RESOLUTION_FUNCTION_RANDOMINT:
		return fromRandomInt(parameters)
	case RESOLUTION_FUNCTION_RANDOMFLOAT:
		return fromRandomFloat(parameters)
	case RESOLUTION_FUNCTION_RANDOMSTRING:
		return fromRandomString(parameters)
	case RESOLUTION_FUNCTION_CONCAT:
		return fromConcat(parameters)
	case RESOLUTION_FUNCTION_ADD:
		return fromAdd(parameters)
	case RESOLUTION_FUNCTION_SUB:
		return fromSub(parameters)
	case RESOLUTION_FUNCTION_MUL:
		return fromMul(parameters)
	case RESOLUTION_FUNCTION_DIV:
		return fromDiv(parameters)
	case RESOLUTION_FUNCTION_FORMATFLOAT:
		return fromFormatFloat(parameters)
	case RESOLUTION_FUNCTION_NOW:
		return fromNow()
	case RESOLUTION_FUNCTION_TODAY:
		return fromToday()
	case RESOLUTION_FUNCTION_TIMESTAMP:
		return fromTimestamp()
	default:
		log.Error().Msgf("Unsupported function [%s], no value is resolved!", functionName)
		return nil
	}
}

// fromUUID generates a new UUID4 string. If an error occurs during UUID generation,
// it logs the error and returns an empty string.
//
// Example: "550e8400-e29b-41d4-a716-446655440000"
func fromUUID() any {

	uuid, error := uuid.NewRandom()
	if error != nil {
		log.Error().Msgf("Impossible to generate value from [uuid] function: %s", error)
		return ""
	}
	return uuid.String()
}

// fromRandomInt generates a random integer between the specified min and max values (inclusive).
// It expects parameters to contain exactly two numeric parameters: min and max. If the parameters
// are invalid or if min is greater than max, it logs an error and returns 0.
//
// Example: parameters = [1, 10] generates a random integer between 1 and 10.
func fromRandomInt(parameters []any) any {

	min, isConversionOk := conversion.AsNumber(parameters[0])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [randomInt] function: first parameter is not a number [%v]", parameters[0])
		return 0
	}

	max, isConversionOk := conversion.AsNumber(parameters[1])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [randomInt] function: second parameter is not a number [%v]", parameters[1])
		return 0
	}

	if min > max {
		log.Error().Msgf("Impossible to generate value from [randomInt] function: first parameter is greater than second parameter [%v > %v]", min, max)
		return 0
	}

	return int64(min + rand.Float64()*(max-min+1))
}

// fromRandomFloat generates a random float between the specified min and max values. It expects
// parameters to contain exactly two numeric parameters: min and max. If the parameters are invalid
// or if min is greater than max, it logs an error and returns 0.0.
//
// Example: parameters = [1.5, 10.5] generates a random float between 1.5 and 10.5 (inclusive).
func fromRandomFloat(parameters []any) any {

	min, isConversionOk := conversion.AsNumber(parameters[0])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [randomFloat] function: first parameter is not a number [%v]", parameters[0])
		return 0
	}

	max, isConversionOk := conversion.AsNumber(parameters[1])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [randomFloat] function: second parameter is not a number [%v]", parameters[1])
		return 0
	}

	if min > max {
		log.Error().Msgf("Impossible to generate value from [randomFloat] function: first parameter is greater than second parameter [%v > %v]", min, max)
		return 0
	}

	return min + rand.Float64()*(max-min)
}

// fromRandomString generates a random alphanumeric string of the specified length. It expects
// parameters to contain exactly one numeric parameter: length. If the parameter is invalid or
// less than or equal to zero, it logs an error and returns an empty string.
//
// Example: parameters = [8] generates a random string of length 8, e.g., "aB3dE9xY".
func fromRandomString(parameters []any) any {

	stringLengthAsFloat, isConversionOk := conversion.AsNumber(parameters[0])
	if !isConversionOk || stringLengthAsFloat <= 0 {
		log.Error().Msgf("Impossible to generate value from [randomString] function: parameter is not a valid string size [%v]", stringLengthAsFloat)
		return ""
	}

	stringLengthAsInt := int(stringLengthAsFloat)
	stringBuffer := make([]byte, stringLengthAsInt)
	charsetLen := len(ALPHANUMERIC_BASIC_CHARSET)
	for i := range stringLengthAsInt {
		stringBuffer[i] = ALPHANUMERIC_BASIC_CHARSET[rand.IntN(charsetLen)]
	}
	return unsafe.String(&stringBuffer[0], stringLengthAsInt)
}

// fromConcat returns the concatenation of two strings. It expects parameters to contain exactly
// two string parameters: prefix and suffix. If the parameters are not parseable as string, it
// logs an error and returns an empty string.
//
// Example: parameters = ["REM", "ORA"] returns "REMORA".
func fromConcat(parameters []any) any {

	prefix, isConversionOk := parameters[0].(string)
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [concat] function: prefix is not a string [%v]", parameters[0])
		return ""
	}

	suffix, isConversionOk := parameters[1].(string)
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [concat] function: suffix is not a string [%v]", parameters[1])
		return ""
	}

	result := prefix + suffix
	log.Trace().Msgf("Executing [concat] function with prefix [%v] and suffix [%v]. Result [%v]", prefix, suffix, result)

	return result
}

// fromAdd returns the sum of two numbers. It expects parameters to contain exactly two numeric
// parameters: addend1 and addend2. If the parameters are invalid or if the result is not a valid
// number, it logs an error and returns 0.
//
// Example: parameters = [5, 3] returns 8.
func fromAdd(parameters []any) any {

	addend1, isConversionOk := conversion.AsNumber(parameters[0])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [add] function: first addend is not a number [%v]", parameters[0])
		return 0
	}

	addend2, isConversionOk := conversion.AsNumber(parameters[1])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [add] function: second addend is not a number [%v]", parameters[1])
		return 0
	}

	result := addend1 + addend2

	if math.IsInf(result, 0) || math.IsNaN(result) {
		log.Error().Msgf("Impossible to generate value from [add] function: result is not a valid number [%v]", result)
		return 0
	}
	log.Trace().Msgf("Executing [add] function with addends [%v] and [%v]. Result [%v]", addend1, addend2, result)

	return result
}

// fromSub returns the difference between two numbers. It expects parameters to contain exactly two
// numeric parameters: minuend and subtrahend. If the parameters are invalid or if the result is not
// a valid number, it logs an error and returns 0.
//
// Example: parameters = [10, 4] returns 6.
func fromSub(parameters []any) any {

	minuend, isConversionOk := conversion.AsNumber(parameters[0])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [sub] function: minuend is not a number [%v]", parameters[0])
		return 0
	}

	subtrahend, isConversionOk := conversion.AsNumber(parameters[1])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [sub] function: subtrahend is not a number [%v]", parameters[1])
		return 0
	}

	result := minuend - subtrahend

	if math.IsInf(result, 0) || math.IsNaN(result) {
		log.Error().Msgf("Impossible to generate value from [sub] function: result is not a valid number [%v]", result)
		return 0
	}
	log.Trace().Msgf("Executing [sub] function with minuend [%v] and subtrahend [%v]. Result [%v]", minuend, subtrahend, result)

	return result
}

// fromMul returns the product of two numbers. It expects parameters to contain exactly two numeric
// parameters: factor1 and factor2. If the parameters are invalid or if the result is not a valid
// number, it logs an error and returns 0.
//
// Example: parameters = [4, 5] returns 20.
func fromMul(parameters []any) any {

	factor1, isConversionOk := conversion.AsNumber(parameters[0])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [mul] function: first factor is not a number [%v]", parameters[0])
		return 0
	}

	factor2, isConversionOk := conversion.AsNumber(parameters[1])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [mul] function: second factor is not a number [%v]", parameters[1])
		return 0
	}

	result := factor1 * factor2

	if math.IsInf(result, 0) || math.IsNaN(result) {
		log.Error().Msgf("Impossible to generate value from [mul] function: result is not a valid number [%v]", result)
		return 0
	}
	log.Trace().Msgf("Executing [mul] function with factors [%v] and [%v]. Result [%v]", factor1, factor2, result)

	return result
}

// fromDiv returns the quotient of two numbers. It expects parameters to contain exactly two numeric
// parameters: numerator and denominator. If the parameters are invalid, if the denominator is zero,
// or if the result is not a valid number, it logs an error and returns 0.
//
// Example: parameters = [20, 4] returns 5.
func fromDiv(parameters []any) any {

	numerator, isConversionOk := conversion.AsNumber(parameters[0])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [div] function: numerator is not a number [%v]", parameters[0])
		return 0
	}

	denominator, isConversionOk := conversion.AsNumber(parameters[1])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [div] function: denominator is not a number [%v]", parameters[1])
		return 0
	}

	if denominator == 0 {
		log.Error().Msgf("Impossible to generate value from [div] function: denominator is zero and division by zero is not allowed")
		return 0
	}

	result := numerator / denominator

	if math.IsInf(result, 0) || math.IsNaN(result) {
		log.Error().Msgf("Impossible to generate value from [div] function: result is not a valid number [%v]", result)
		return 0
	}
	log.Trace().Msgf("Executing [div] function with numerator [%v] and denominator [%v]. Result [%v]", numerator, denominator, result)

	return result
}

// fromFormatFloat formats a float to the specified number of decimal places. It expects parameters to
// contain exactly two parameters: number (float) and precision (non-negative integer). If the parameters
// are invalid or if the result is not a valid formatted float, it logs an error and returns an empty string.
//
// Example: parameters = [3.14159, 2] returns 3.14.
func fromFormatFloat(parameters []any) any {

	number, isConversionOk := conversion.AsNumber(parameters[0])
	if !isConversionOk {
		log.Error().Msgf("Impossible to generate value from [formatFloat] function: first parameter is not a number [%v]", parameters[0])
		return ""
	}

	precision, isConversionOk := conversion.AsNumber(parameters[1])
	if !isConversionOk || precision < 0 || precision != math.Trunc(precision) {
		log.Error().Msgf("Impossible to generate value from [formatFloat] function: second parameter is not a valid precision number [%v]", parameters[1])
		return ""
	}

	result := strconv.FormatFloat(number, 'f', int(precision), 64)
	if result == "" {
		log.Error().Msgf("Impossible to generate value from [formatFloat] function: result is not a valid formatted float [%v]", result)
		return ""
	}
	log.Trace().Msgf("Executing [formatFloat] function with number [%v] and precision [%v]. Result [%v]", number, precision, result)

	return result
}

// fromNow returns the current date and time as a formatted string "YYYY-MM-DD HH:MM:SS". If an error
// occurs during formatting, it logs the error and returns an empty string.
//
// Example: "2023-10-05 14:30:14"
func fromNow() any {

	result := time.Now().Format("2006-01-02 15:04:05")
	if result == "" {
		log.Error().Msgf("Impossible to generate value from [now] function: result is not a valid datetime string [%v]", result)
		return ""
	}
	log.Trace().Msgf("Executing [now] function. Result [%v]", result)

	return result
}

// fromToday returns the current date as a formatted string "YYYY-MM-DD". If an error occurs during
// formatting, it logs the error and returns an empty string.
//
// Example: "2023-10-05"
func fromToday() any {

	result := time.Now().Format("2006-01-02")
	if result == "" {
		log.Error().Msgf("Impossible to generate value from [today] function: result is not a valid date string [%v]", result)
		return ""
	}
	log.Trace().Msgf("Executing [today] function. Result [%v]", result)

	return result
}

// fromTimestamp returns the current Unix timestamp as an integer. If an error occurs during timestamp
// generation, it logs the error and returns 0.
//
// Example: 1735729200000
func fromTimestamp() any {

	result := time.Now().Unix()
	if result <= 0 {
		log.Error().Msgf("Impossible to generate value from [timestamp] function: result is not a valid timestamp [%v]", result)
		return 0
	}
	log.Trace().Msgf("Executing [timestamp] function. Result [%v]", result)

	return result
}
