package operator

import (
	"remora/pkg/conversion"
	"remora/pkg/customerror"

	"github.com/rs/zerolog/log"
)

// IsNull evaluates whether the provided value is nil.
//
// It returns true if the input is nil, otherwise false.
func IsNull(evaluatedValue any) (bool, error) {

	isCompliant := evaluatedValue == nil
	log.Trace().Msgf("Evaluating [is_null]: value [%v] is null? [%v]", evaluatedValue, isCompliant)
	return isCompliant, nil
}

// IsEmpty checks whether the provided value is an empty string.
//
// It returns true if the input is a string with zero length.
// It returns an error if evaluatedValue is not a string value.
func IsEmpty(evaluatedValue any) (bool, error) {

	value, ok := evaluatedValue.(string)
	if !ok {
		log.Trace().Msgf("Evaluating [is_empty]: value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	isCompliant := len(value) == 0
	log.Trace().Msgf("Evaluating [is_empty]: value [%v] is empty? [%v]", evaluatedValue, isCompliant)
	return isCompliant, nil
}

// IsNotNull checks whether the provided value is not nil.
//
// It returns true if the input value is non-nil, false otherwise.
func IsNotNull(evaluatedValue any) (bool, error) {

	isCompliant := evaluatedValue != nil
	log.Trace().Msgf("Evaluating [is_not_null]: value [%v] is not null? [%v]", evaluatedValue, isCompliant)
	return isCompliant, nil
}

// IsNotEmpty checks whether the provided value is a non-empty string.
//
// It returns true if evaluatedValue is a string with length greater than zero.
// It returns an error if evaluatedValue is not a string value.
func IsNotEmpty(evaluatedValue any) (bool, error) {

	value, ok := evaluatedValue.(string)
	if !ok {
		log.Trace().Msgf("Evaluating [is_not_empty]: value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	isCompliant := len(value) > 0
	log.Trace().Msgf("Evaluating [is_not_empty]: value [%v] is not empty? [%v]", evaluatedValue, isCompliant)
	return isCompliant, nil
}

// IsTrue evaluates whether the provided value represents a boolean true.
//
// It returns true if evaluatedValue is a boolean with true value.
// It returns an error if evaluatedValue is not a boolean value.
func IsTrue(evaluatedValue any) (bool, error) {

	value, isOk := conversion.AsBoolean(evaluatedValue)
	if !isOk {
		log.Trace().Msgf("Evaluating [is_true]: value [%v] is not a boolean", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidBooleanConversion, evaluatedValue)
	}
	log.Trace().Msgf("Evaluating [is_true]: value [%v] is true? [%v]", evaluatedValue, value)
	return value, nil
}

// IsFalse evaluates whether the provided value represents a boolean false.
//
// It returns true if evaluatedValue is a boolean with false value.
// It returns an error if evaluatedValue is not a boolean value.
func IsFalse(evaluatedValue any) (bool, error) {

	value, isOk := conversion.AsBoolean(evaluatedValue)
	if !isOk {
		log.Trace().Msgf("Evaluating [is_true]: value [%v] is not a boolean", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidBooleanConversion, evaluatedValue)
	}
	log.Trace().Msgf("Evaluating [is_false]: value [%v] is false? [%v]", evaluatedValue, !value)
	return !value, nil
}
