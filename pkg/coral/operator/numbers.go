package operator

import (
	"remora/pkg/conversion"
	"remora/pkg/customerror"

	"github.com/rs/zerolog/log"
)

// NumberEqualsTo checks whether evaluatedValue equals comparationValue as number.
//
// Returns true only if both values are number and exactly equal.
// It returns an error if evaluatedValue or comparationValue is not a number value.
func NumberEqualsTo(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldNumber := conversion.AsNumber(evaluatedValue)
	if !isFieldNumber {
		log.Trace().Msgf("Evaluating [eq_n]: evaluated value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	clauseValue, isClauseNumber := conversion.AsNumber(comparationValue)
	if !isClauseNumber {
		log.Trace().Msgf("Evaluating [eq_n]: clause value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	isCompliant := isFieldNumber && isClauseNumber && value == clauseValue
	log.Trace().Msgf("Evaluating [eq_n]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, clauseValue, clauseValue, isCompliant)
	return isCompliant, nil
}

// NumberNotEqualsTo checks whether evaluatedValue differs comparationValue as number.
//
// Returns true only if both values are number and not equal.
// It returns an error if evaluatedValue or comparationValue is not a number value.
func NumberNotEqualsTo(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldNumber := conversion.AsNumber(evaluatedValue)
	if !isFieldNumber {
		log.Trace().Msgf("Evaluating [neq_n]: evaluated value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	clauseValue, isClauseNumber := conversion.AsNumber(comparationValue)
	if !isClauseNumber {
		log.Trace().Msgf("Evaluating [neq_n]: clause value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	isCompliant := isFieldNumber && isClauseNumber && value != clauseValue
	log.Trace().Msgf("Evaluating [neq_n]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, clauseValue, clauseValue, isCompliant)
	return isCompliant, nil
}

// NumberGreaterThan checks whether evaluatedValue is greater than comparationValue as number.
//
// Returns true only if both values are number and comparationValue is greater than evaluatedValue.
// It returns an error if evaluatedValue or comparationValue is not a number value.
func NumberGreaterThan(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldNumber := conversion.AsNumber(evaluatedValue)
	if !isFieldNumber {
		log.Trace().Msgf("Evaluating [gt_n]: evaluated value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	clauseValue, isClauseNumber := conversion.AsNumber(comparationValue)
	if !isClauseNumber {
		log.Trace().Msgf("Evaluating [gt_n]: clause value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	isCompliant := value > clauseValue
	log.Trace().Msgf("Evaluating [gt_n]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, clauseValue, clauseValue, isCompliant)
	return isCompliant, nil
}

// NumberGreaterThanOrEqualsTo checks whether evaluatedValue is greater or equals than comparationValue as number.
//
// Returns true only if both values are number and comparationValue is greater or equals than evaluatedValue.
// It returns an error if evaluatedValue or comparationValue is not a number value.
func NumberGreaterThanOrEqualsTo(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldNumber := conversion.AsNumber(evaluatedValue)
	if !isFieldNumber {
		log.Trace().Msgf("Evaluating [gte_n]: evaluated value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	clauseValue, isClauseNumber := conversion.AsNumber(comparationValue)
	if !isClauseNumber {
		log.Trace().Msgf("Evaluating [gte_n]: clause value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	isCompliant := value >= clauseValue
	log.Trace().Msgf("Evaluating [gte_n]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, clauseValue, clauseValue, isCompliant)
	return isCompliant, nil
}

// NumberLowerThan checks whether evaluatedValue is lower than comparationValue as number.
//
// Returns true only if both values are number and comparationValue is lower than evaluatedValue.
// It returns an error if evaluatedValue or comparationValue is not a number value.
func NumberLowerThan(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldNumber := conversion.AsNumber(evaluatedValue)
	if !isFieldNumber {
		log.Trace().Msgf("Evaluating [lt_n]: evaluated value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	clauseValue, isClauseNumber := conversion.AsNumber(comparationValue)
	if !isClauseNumber {
		log.Trace().Msgf("Evaluating [lt_n]: clause value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	isCompliant := value < clauseValue
	log.Trace().Msgf("Evaluating [lt_n]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, clauseValue, clauseValue, isCompliant)
	return isCompliant, nil
}

// NumberLowerThanOrEqualsTo checks whether evaluatedValue is lower or equals than comparationValue as number.
//
// Returns true only if both values are number and comparationValue is lower or equals than evaluatedValue.
// It returns an error if evaluatedValue or comparationValue is not a number value.
func NumberLowerThanOrEqualsTo(comparationValue any, evaluatedValue any) (bool, error) {

	value, isFieldNumber := conversion.AsNumber(evaluatedValue)
	if !isFieldNumber {
		log.Trace().Msgf("Evaluating [lte_n]: evaluated value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	clauseValue, isClauseNumber := conversion.AsNumber(comparationValue)
	if !isClauseNumber {
		log.Trace().Msgf("Evaluating [lte_n]: clause value [%v] is not a number", value)
		return false, customerror.NewError(customerror.ErrorCoralOperatorInvalidNumericConversion, evaluatedValue)
	}

	isCompliant := value <= clauseValue
	log.Trace().Msgf("Evaluating [lte_n]: value [%v] (type [%T]) and clause value [%v] (type [%T]) are compliant? [%v]", value, value, clauseValue, clauseValue, isCompliant)
	return isCompliant, nil
}
