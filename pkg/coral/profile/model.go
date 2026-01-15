package profile

import (
	"remora/pkg/coral/spec/operator"
)

// Profile contains the set of behaviors associated to a Minnow element that
// defines the possible effects that a Minnow can reproduce regarding prompted
// input data.
type Profile struct {

	// Behaviors provides a list of [Behavior] elements that compose the
	// Minnow behavioral profile
	Behaviors []Behavior
}

// Behavior represents a single Minnow's behavioral attitude under certain conditions,
// expliciting the required [Condition] that produce specific [Effect].
type Behavior struct {

	// Condition provides the set of clauses, chained together in order to define
	// the steps of a validation flow
	Condition Condition

	// Effect provides an action or a content to return if and only if all the clauses
	// in the condition have been satisfied
	Effect Effect
}

// Condition represents one of the conditions that are chained together to form a
// complex validation pattern for Minnow behavior
type Condition struct {

	// Terms provides the conditions on nested level regarding the present condition.
	// This could be empty if no nested condition is present
	Terms []Condition

	// Type provides the logical conjunction that links this condition to the previous
	// one in the validation chain.
	Type operator.ExpressionType

	// Clause provides the composition of the condition regarding the structure of the
	// controls to be validated. This could be empty if nested condition is present.
	Clause *Clause
}

// Clause represents the atomic conditional check that can be part of a check action
// in the validation chain. It can contain a binary or unary operator, so it may or
// may not define a specific value.
type Clause struct {

	// Field provides the name of the field that can be found in the input data during
	// the validation process and so it is used in the check step.
	Field string

	// Operator provides the type of operator applied during the validation process of
	// the current clause. It can be binary or unary.
	Operator operator.Operator

	// Value provide the value on which the operator is applied during the validation
	// process of the current clause. It can be a static value or a placeholder value,
	// used for data resolution. If the operator is unary, it must be nil.
	Value any
}

// Effect represents the consequential effect of Minnow behavior if all condition in the
// chain validation are met. It can define an action, where a complex execution is
// performed, or a value return.
type Effect struct {

	// Action provides the complex operation to be executed if an effect is reached
	Action *Action

	// Action provides the static and dynamic value to be returned if an effect is reached
	Return *Return
}

// Action represents the action that can be executed after the chain validation as effect
// to certain Minnow's behavioral pattern
type Action struct {

	// TODO: currently not used, just a placeholder
	Code string
}

// Return represents the value to return after the chain validation as effect to certain
// Minnow's behavioral pattern. It can contain static values (that are returned as-is)
// and dynamic values (that needs data resolution).
type Return struct {

	// Value provides the static or dynamic element to returns in behavior effect
	Value string
}

// EffectResult represents the resulting effect generated following completion of the
// validation chain on certain behavior, only when all checks have been successful (or,
// alternatively, are results of the default behavior).
// The resulting effect is represented to client as an HTTP response.
type EffectResult struct {

	// ContentType provides the value of the "Content-Type" HTTP header to be returned
	// in final response.
	ContentType string

	// RawBody provides a raw string that contains the explicit body to be returned in
	// final response.
	RawBody string

	// Headers provides the value of the HTTP headers to be returned in final response.
	Headers map[string]any

	// MappableBody provides a map that contains the key-value content to be used for generate
	// a body in custom content type to be returned in final response.
	MappableBody map[string]any

	// StatusCode provides the HTTP status code to be returned in final response.
	StatusCode int

	// IsValid provides the logical value taken from the validation chain.
	// If the validation chain has been successfully completed, the value is true,
	// otherwise it is false.
	IsValid bool
}
