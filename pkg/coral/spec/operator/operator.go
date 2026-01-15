package operator

// ExpressionType represents the type of expression joined by a [Condition].
// This is an alias for string to improve code readability and set predefined constants.
// The accepted values are: and, or, clause.
type ExpressionType string

const (

	// EXPR_TYPE_AND define an 'AND' conjunction
	EXPR_TYPE_AND ExpressionType = "and"

	// EXPR_TYPE_AND define an 'OR' conjunction
	EXPR_TYPE_OR ExpressionType = "or"

	// EXPR_TYPE_AND define the presence of nested clauses in same condition
	EXPR_TYPE_CLAUSE ExpressionType = "clause"
)

// Operator represents the operator to be used for validation on behavior analysis.
// This is an alias for string to improve code readability and set predefined constants.
type Operator string

const (
	// -- Unary operators

	// UNARY_ISNULL defines the enumerative value of unary operator for
	// '<field> is null' condition
	UNARY_ISNULL Operator = "is_null"

	// UNARY_ISEMPTY defines the enumerative value of unary operator for
	// '<field> is empty' condition
	UNARY_ISEMPTY Operator = "is_empty"

	// UNARY_ISNOTNULL defines the enumerative value of unary operator for
	// '<field> is not null' condition
	UNARY_ISNOTNULL Operator = "is_not_null"

	// UNARY_ISNOTEMPTY defines the enumerative value of unary operator for
	// '<field> is not null' condition
	UNARY_ISNOTEMPTY Operator = "is_not_empty"

	// UNARY_ISTRUE defines the enumerative value of unary operator for
	// '<field> is true' condition
	UNARY_ISTRUE Operator = "is_true"

	// UNARY_ISFALSE defines the enumerative value of unary operator for
	// '<field> is false' condition
	UNARY_ISFALSE Operator = "is_false"

	// -- Binary operators

	// BINARY_EQUALS_STRING defines the enumerative value of binary operator
	// for '<field> == <string>' condition
	BINARY_EQUALS_STRING Operator = "eq_s"

	// BINARY_EQUALS_CASEINSENSITIVE_STRING defines the enumerative value of
	// binary operator for case-insensitive '<field> == <string>' condition
	BINARY_EQUALS_CASEINSENSITIVE_STRING Operator = "eqci_s"

	// BINARY_EQUALS_NUMBER defines the enumerative value of binary operator
	// for '<field> == <number>' condition
	BINARY_EQUALS_NUMBER Operator = "eq_n"

	// BINARY_NOTEQUALS_STRING defines the enumerative value of binary operator
	// for '<field> != <string>' condition
	BINARY_NOTEQUALS_STRING Operator = "neq_s"

	// BINARY_NOTEQUALS_CASEINSENSITIVE_STRING defines the enumerative value of
	// binary operator for case-insensitive '<field> != <string>' condition
	BINARY_NOTEQUALS_CASEINSENSITIVE_STRING Operator = "neqci_s"

	// BINARY_NOTEQUALS_NUMBER defines the enumerative value of binary operator
	// for '<field> != <number>' condition
	BINARY_NOTEQUALS_NUMBER Operator = "neq_n"

	// BINARY_GREATERTHAN_NUMBER defines the enumerative value of binary operator
	// for '<field> > <number>' condition
	BINARY_GREATERTHAN_NUMBER Operator = "gt_n"

	// BINARY_GREATEREQUALSTHAN_NUMBER defines the enumerative value of binary operator
	// for '<field> >= <number>' condition
	BINARY_GREATEREQUALSTHAN_NUMBER Operator = "gte_n"

	// BINARY_LOWERTHAN_NUMBER defines the enumerative value of binary operator for
	// '<field> < <number>' condition
	BINARY_LOWERTHAN_NUMBER Operator = "lt_n"

	// BINARY_LOWEREQUALSTHAN_NUMBER defines the enumerative value of binary operator
	// for '<field> <= <number>' condition
	BINARY_LOWEREQUALSTHAN_NUMBER Operator = "lte_n"

	// BINARY_STARTSWITH defines the enumerative value of binary operator for
	// '<field> startswith <string>' condition
	BINARY_STARTSWITH Operator = "startswith"

	// BINARY_ENDSWITH defines the enumerative value of binary operator for
	// '<field> endswith <string>' condition
	BINARY_ENDSWITH Operator = "endswith"

	// BINARY_CONTAINS_SUBSTRING defines the enumerative value of binary operator for
	// '<field> contains <string>' condition
	BINARY_CONTAINS_SUBSTRING Operator = "contains_s"

	// BINARY_REGEX defines the enumerative value of binary operator for
	// '<field> is compliant to <regex>' condition
	BINARY_REGEX Operator = "regex"

	// SPECIAL_DEFAULT defines the enumerative value of special operator that provides
	// the "always true" condition in OTHERWISE/DEFAULT cases
	SPECIAL_DEFAULT Operator = "as_default"
)

// UnaryOperators defines the list of unary operators, as Operator type
var UnaryOperators = []Operator{
	UNARY_ISNULL, UNARY_ISEMPTY, UNARY_ISNOTNULL,
	UNARY_ISNOTEMPTY, UNARY_ISTRUE, UNARY_ISFALSE,
}

// BinaryOperators defines the list of binary operators, as Operator type
var BinaryOperators = []Operator{
	BINARY_EQUALS_STRING, BINARY_EQUALS_NUMBER,
	BINARY_NOTEQUALS_STRING, BINARY_NOTEQUALS_NUMBER,
	BINARY_GREATERTHAN_NUMBER, BINARY_GREATEREQUALSTHAN_NUMBER,
	BINARY_LOWERTHAN_NUMBER, BINARY_LOWEREQUALSTHAN_NUMBER,
	BINARY_STARTSWITH, BINARY_ENDSWITH, BINARY_CONTAINS_SUBSTRING, BINARY_REGEX,
}
