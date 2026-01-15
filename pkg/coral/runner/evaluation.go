package runner

import (
	"remora/pkg/conversion/maps"
	ops "remora/pkg/coral/operator"
	"remora/pkg/coral/profile"
	"remora/pkg/coral/spec/operator"
	"remora/pkg/customerror"
	"slices"
)

// isConditionSatisfied evaluates whether the given condition is satisfied by the provided input map.
//
// The evaluation follows the logical structure of the condition:
//   - CLAUSE conditions are evaluated directly against the input map
//   - AND conditions are satisfied only if all nested conditions are satisfied
//   - OR conditions are satisfied if at least one nested condition is satisfied
//
// Evaluation is performed recursively and stops early whenever the final outcome
// can be determined. Any error raised during clause evaluation is propagated to the caller.
func isConditionSatisfied(condition profile.Condition, input map[string]any) (bool, error) {

	switch condition.Type {

	// Atomic clause reached: check straightforwardly if is satisfied
	case operator.EXPR_TYPE_CLAUSE:
		isCompliant, err := isClauseSatisfied(*condition.Clause, input)
		return isCompliant, err

	// AND clause reached: check if EACH ONE is satisfied
	case operator.EXPR_TYPE_AND:
		for _, term := range condition.Terms {
			isCompliant, err := isConditionSatisfied(term, input)
			if err != nil {
				return false, err
			}
			if !isCompliant {
				return false, nil
			}
		}
		return true, nil

	// OR clause reached: check if AT LEAST ONE is satisfied
	case operator.EXPR_TYPE_OR:
		for _, term := range condition.Terms {
			isCompliant, err := isConditionSatisfied(term, input)
			if err != nil {
				return false, err
			}
			if isCompliant {
				return true, nil
			}
		}
		return false, nil

	default:
		return false, customerror.NewError(customerror.ErrorCoralSyntaxTreeUnmappedClause, condition.Type)
	}
}

// isClauseSatisfied evaluates whether the given clause is satisfied by the provided input map.
//
// Clause evaluation depends on the operator type:
//   - binary operators compare a field against a value
//   - unary operators validate the presence or state of a field without an explicit value
//   - special operators define engine-level behaviors and may bypass standard evaluation
//
// The function returns an error only if the underlying operator evaluation fails.
func isClauseSatisfied(clause profile.Clause, input map[string]any) (bool, error) {

	op := clause.Operator

	if op == operator.SPECIAL_DEFAULT {
		return true, nil
	}

	if slices.Contains(operator.BinaryOperators, op) {
		return isBinaryOperatorSatisfied(clause, input)
	}

	if slices.Contains(operator.UnaryOperators, op) {
		return isUnaryOperatorSatisfied(clause, input)
	}

	return false, nil
}

// isBinaryOperatorSatisfied evaluates whether the given clause is satisfied by applying
// a binary operator to a field value extracted from the input map.
//
// Binary operators compare the field value against a clause value, which may be defined
// statically or resolved dynamically at runtime. The specific comparison logic depends
// on the operator type (equality, ordering, string matching or pattern matching).
//
// An error is returned if the operator is not supported or if the evaluation fails.
func isBinaryOperatorSatisfied(clause profile.Clause, input map[string]any) (bool, error) {

	fieldValue := maps.ExtractFieldValue(input, clause.Field)
	op := clause.Operator

	switch op {

	case operator.BINARY_EQUALS_STRING:
		return ops.StringEqualsTo(clause.Value, fieldValue)
	case operator.BINARY_EQUALS_CASEINSENSITIVE_STRING:
		return ops.StringCaseInsensitiveEqualsTo(clause.Value, fieldValue)
	case operator.BINARY_EQUALS_NUMBER:
		return ops.NumberEqualsTo(clause.Value, fieldValue)
	case operator.BINARY_NOTEQUALS_STRING:
		return ops.StringNotEqualsTo(clause.Value, fieldValue)
	case operator.BINARY_NOTEQUALS_CASEINSENSITIVE_STRING:
		return ops.StringCaseInsensitiveNotEqualsTo(clause.Value, fieldValue)
	case operator.BINARY_NOTEQUALS_NUMBER:
		return ops.NumberNotEqualsTo(clause.Value, fieldValue)
	case operator.BINARY_GREATERTHAN_NUMBER:
		return ops.NumberGreaterThan(clause.Value, fieldValue)
	case operator.BINARY_GREATEREQUALSTHAN_NUMBER:
		return ops.NumberGreaterThanOrEqualsTo(clause.Value, fieldValue)
	case operator.BINARY_LOWERTHAN_NUMBER:
		return ops.NumberLowerThan(clause.Value, fieldValue)
	case operator.BINARY_LOWEREQUALSTHAN_NUMBER:
		return ops.NumberLowerThanOrEqualsTo(clause.Value, fieldValue)
	case operator.BINARY_STARTSWITH:
		return ops.StringStartsWith(clause.Value, fieldValue)
	case operator.BINARY_ENDSWITH:
		return ops.StringEndsWith(clause.Value, fieldValue)
	case operator.BINARY_CONTAINS_SUBSTRING:
		return ops.StringContainsSubstring(clause.Value, fieldValue)
	case operator.BINARY_REGEX:
		return ops.StringCompliantToRegex(clause.Value, fieldValue)
	default:
		return false, customerror.NewError(customerror.ErrorCoralSyntaxTreeOperatorNotImplemented, op)
	}
}

// isUnaryOperatorSatisfied evaluates whether the given clause is satisfied by applying
// a unary operator to a field value extracted from the input map.
//
// Unary operators perform validations that depend only on the presence, state, or
// boolean interpretation of the field value, without comparing it to an explicit
// clause value.
//
// An error is returned if the operator is not supported or if the evaluation fails.
func isUnaryOperatorSatisfied(clause profile.Clause, input map[string]any) (bool, error) {

	fieldValue := maps.ExtractFieldValue(input, clause.Field)
	op := clause.Operator

	switch op {

	case operator.UNARY_ISNULL:
		return ops.IsNull(fieldValue)
	case operator.UNARY_ISEMPTY:
		return ops.IsEmpty(fieldValue)
	case operator.UNARY_ISNOTNULL:
		return ops.IsNotNull(fieldValue)
	case operator.UNARY_ISNOTEMPTY:
		return ops.IsNotEmpty(fieldValue)
	case operator.UNARY_ISTRUE:
		return ops.IsTrue(fieldValue)
	case operator.UNARY_ISFALSE:
		return ops.IsFalse(fieldValue)
	default:
		return false, customerror.NewError(customerror.ErrorCoralSyntaxTreeOperatorNotImplemented, op)
	}
}
