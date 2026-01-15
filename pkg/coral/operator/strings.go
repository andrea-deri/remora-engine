package operator

import (
	"regexp"
	"remora/pkg/customerror"
	"strings"

	"github.com/rs/zerolog/log"
)

// StringEqualsTo checks whether evaluatedValue equals comparationValue as strings.
//
// Returns true only if both values are strings and exactly equal.
// It returns an error if evaluatedValue or comparationValue is not a string value.
func StringEqualsTo(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldString := evaluatedValue.(string)
	if !isFieldString {
		log.Trace().Msgf("Evaluating [eq_s]: evaluated value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	clauseValue, isClauseString := comparationValue.(string)
	if !isClauseString {
		log.Trace().Msgf("Evaluating [eq_s]: clause value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	isCompliant := value == clauseValue
	log.Trace().Msgf("Evaluating [eq_s]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, comparationValue, comparationValue, isCompliant)
	return isCompliant, nil
}

// StringCaseInsensitiveEqualsTo checks whether evaluatedValue equals comparationValue as strings
// in case-insensitive comparation.
//
// Returns true only if both values are strings and equal, ignoring word cases.
// It returns an error if evaluatedValue or comparationValue is not a string value.
func StringCaseInsensitiveEqualsTo(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldString := evaluatedValue.(string)
	if !isFieldString {
		log.Trace().Msgf("Evaluating [eqci_s]: evaluated value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	clauseValue, isClauseString := comparationValue.(string)
	if !isClauseString {
		log.Trace().Msgf("Evaluating [eqci_s]: clause value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	isCompliant := strings.EqualFold(value, clauseValue)
	log.Trace().Msgf("Evaluating [eqci_s]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, comparationValue, comparationValue, isCompliant)
	return isCompliant, nil
}

// StringNotEqualsTo checks whether evaluatedValue differs comparationValue as strings.
//
// Returns true only if both values are strings and not equal.
// It returns an error if evaluatedValue or comparationValue is not a string value.
func StringNotEqualsTo(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldString := evaluatedValue.(string)
	if !isFieldString {
		log.Trace().Msgf("Evaluating [neq_s]: evaluated value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	clauseValue, isClauseString := comparationValue.(string)
	if !isClauseString {
		log.Trace().Msgf("Evaluating [eqci_s]: clause value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	isCompliant := value != clauseValue
	log.Trace().Msgf("Evaluating [neq_s]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, comparationValue, comparationValue, isCompliant)
	return isCompliant, nil
}

// StringCaseInsensitiveNotEqualsTo checks whether evaluatedValue differs comparationValue as strings
// in case-insensitive comparation.
//
// Returns true only if both values are strings and not equal, ignoring word cases.
// It returns an error if evaluatedValue or comparationValue is not a string value.
func StringCaseInsensitiveNotEqualsTo(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldString := evaluatedValue.(string)
	if !isFieldString {
		log.Trace().Msgf("Evaluating [neqci_s]: evaluated value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	clauseValue, isClauseString := comparationValue.(string)
	if !isClauseString {
		log.Trace().Msgf("Evaluating [neqci_s]: clause value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	isCompliant := !strings.EqualFold(value, clauseValue)
	log.Trace().Msgf("Evaluating [neqci_s]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, comparationValue, comparationValue, isCompliant)
	return isCompliant, nil
}

// StringStartsWith checks whether evaluatedValue starts with comparationValue as prefix.
//
// Returns true only if both values are strings and prefix condition is met.
// It returns an error if evaluatedValue or comparationValue is not a string value.
func StringStartsWith(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldString := evaluatedValue.(string)
	if !isFieldString {
		log.Trace().Msgf("Evaluating [startswith]: evaluated value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	clauseValue, isClauseString := comparationValue.(string)
	if !isClauseString {
		log.Trace().Msgf("Evaluating [startswith]: clause value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	isCompliant := strings.HasPrefix(value, clauseValue)
	log.Trace().Msgf("Evaluating [startswith]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, comparationValue, clauseValue, clauseValue, isCompliant)
	return isCompliant, nil
}

// StringEndsWith checks whether evaluatedValue ends with comparationValue as suffix.
//
// Returns true only if both values are strings and suffix condition is met.
// It returns an error if evaluatedValue or comparationValue is not a string value.
func StringEndsWith(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldString := evaluatedValue.(string)
	if !isFieldString {
		log.Trace().Msgf("Evaluating [endswith]: evaluated value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	clauseValue, isClauseString := comparationValue.(string)
	if !isClauseString {
		log.Trace().Msgf("Evaluating [endswith]: clause value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	isCompliant := strings.HasSuffix(value, clauseValue)
	log.Trace().Msgf("Evaluating [endswith]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, comparationValue, clauseValue, clauseValue, isCompliant)
	return isCompliant, nil
}

// StringContainsSubstring checks whether evaluatedValue contains comparationValue as substring.
//
// Returns true only if both values are strings and substring condition is met.
// It returns an error if evaluatedValue or comparationValue is not a string value.
func StringContainsSubstring(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldString := evaluatedValue.(string)
	if !isFieldString {
		log.Trace().Msgf("Evaluating [contains_s]: evaluated value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	clauseValue, isClauseString := comparationValue.(string)
	if !isClauseString {
		log.Trace().Msgf("Evaluating [contains_s]: clause value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	isCompliant := strings.Contains(value, clauseValue)
	log.Trace().Msgf("Evaluating [contains_s]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, comparationValue, clauseValue, clauseValue, isCompliant)
	return isCompliant, nil
}

// StringCompliantToRegex checks whether evaluatedValue match comparationValue as compiled regex.
//
// Returns true only if evaluatedValue is a string value, comparationValue is a regex value and
// substring condition is met.
// It returns an error if evaluatedValue is not a string value or comparationValue is not a regex value.
func StringCompliantToRegex(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldString := evaluatedValue.(string)
	if !isFieldString {
		log.Trace().Msgf("Evaluating [regex]: evaluated value [%v] is not a string", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidStringConversion, evaluatedValue)
	}

	clauseValue, isCompiledRegex := comparationValue.(*regexp.Regexp)
	if !isCompiledRegex {
		log.Trace().Msgf("Evaluating [regex]: clause value [%v] is not a regex", value)
		return false, customerror.NewError(customerror.ErrorCoralSyntaxTreeRegexNotApplicable, evaluatedValue)
	}

	isCompliant := clauseValue.MatchString(value)
	log.Trace().Msgf("Evaluating [regex]: value [%v] (type [Regexp]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, comparationValue, isCompliant)
	return isCompliant, nil
}
