package customerror

import (
	"fmt"
)

const (

	// REMORA_ERROR_STRUCTURE represents the formatting structure for error codes in REMORA
	REMORA_ERROR_STRUCTURE = "RMR-%s"
)

// RemoraError represents a structured error with a code and message.
type RemoraError struct {

	// Code defines the specific numerical error code, that will
	// appears in the form defined by REMORA_ERROR_STRUCTURE constant.
	Code RemoraErrorCode

	// Message defines the error message, that can contains format verbs
	// to be replaced by arguments passed during error creation.
	Message string
}

// NewError creates a new error with the given RemoraError and arguments.
func NewError(kind RemoraError, args ...any) error {

	msg := fmt.Sprintf(kind.Message, args...)
	return fmt.Errorf("[%s] - %s", fmt.Sprintf(REMORA_ERROR_STRUCTURE, kind.Code), msg)
}

// ErrorGZipNotDecodable is returned when a GZIP string cannot be decoded.
var ErrorGZipNotDecodable = RemoraError{
	Code:    CODE_HARBOR_GZIP_NOT_DECODABLE,
	Message: "Impossible to decode input string from GZIP format. %v",
}

// ErrorTidalSearchProfileNotFound is returned when no valid ruleset
// is found in profile map.
var ErrorTidalSearchProfileNotFound = RemoraError{
	Code:    CODE_TIDAL_SEARCH_PROFILE_NOT_FOUND,
	Message: "No valid Profile found with ID [%s]",
}

// ErrorCoralSyntaxTreeUnmappedClause is returned when a clause is
// unmapped in syntax tree.
var ErrorCoralSyntaxTreeUnmappedClause = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_UNMAPPED_CLAUSE,
	Message: "Unmapped clause [%s], impossible to evaluate",
}

// ErrorCoralSyntaxTreeOperatorNotImplemented is returned when an operator
// is not implemented in syntax tree.
var ErrorCoralSyntaxTreeOperatorNotImplemented = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_OPERATOR_NOT_IMPLEMENTED,
	Message: "Operator [%s] not yet implemented",
}

// ErrorCoralSyntaxTreeFunctionalityNotImplemented is returned when a
// functionality is not implemented in syntax tree.
var ErrorCoralSyntaxTreeFunctionalityNotImplemented = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_FUNCTIONALITY_NOT_IMPLEMENTED,
	Message: "Functionality [%s] not yet implemented",
}

// ErrorCoralSyntaxTreeRegexNotApplicable is returned when regex is
// not applicable in syntax tree.
var ErrorCoralSyntaxTreeRegexNotApplicable = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_REGEX_NOT_APPLICABLE,
	Message: "Regex not applicable on non-expression value [%s]",
}

// ErrorCoralSyntaxTreeMalformedLexeme is returned when a lexeme is malformed
// in syntax tree.
var ErrorCoralSyntaxTreeMalformedLexeme = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_MALFORMED_LEXEME,
	Message: "Generated malformed lexeme with index [%d]: invalid token at start [%s] or end [%s]",
}

// ErrorCoralSyntaxTreeInvalidStartingTokenLexeme is returned when a lexeme
// starts with an invalid token in syntax tree.
var ErrorCoralSyntaxTreeInvalidStartingTokenLexeme = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_INVALID_STARTING_TOKEN_LEXEME,
	Message: "The lexeme did not starts with [IF], [OTHERWISE] or [DEFAULT] token, found [%s] token (%d:%d)",
}

// ErrorCoralSyntaxTreeInvalidEndingTokenLexeme is returned when a lexeme
// ends with an invalid token in syntax tree.
var ErrorCoralSyntaxTreeInvalidEndingTokenLexeme = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_INVALID_ENDING_TOKEN_LEXEME,
	Message: "The lexeme did not ends with [END] token, found [%s] token (%d:%d)",
}

// ErrorCoralSyntaxTreeMissingRightParenthesys is returned when a lexeme
// contains a left parenthesys but not a right parenthesys in syntax tree.
var ErrorCoralSyntaxTreeMissingRightParenthesys = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_MISSING_RIGHT_PARENTHESYS,
	Message: "The lexeme contains a left parenthesys but not a right parenthesys (%d:%d)",
}

// ErrorCoralSyntaxTreeMissingLiteralToken is returned when a lexeme does
// not contain a variable literal token in syntax tree.
var ErrorCoralSyntaxTreeMissingLiteralToken = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_MISSING_LITERAL_TOKEN,
	Message: "The lexeme did not contains the variable literal token, found [%s] token (%d:%d)",
}

// ErrorCoralSyntaxTreeInvalidOperatorToken is returned when a lexeme does
// not contain a valid operator token in syntax tree.
var ErrorCoralSyntaxTreeInvalidOperatorToken = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_INVALID_OPERATOR_TOKEN,
	Message: "The lexeme did not contains a valid operator token, found [%s] token (%d:%d)",
}

// ErrorCoralSyntaxTreeInvalidRightOperatorToken is returned when a lexeme
// does not contain a valid right-operator token in syntax tree.
var ErrorCoralSyntaxTreeInvalidRightOperatorToken = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_INVALID_RIGHT_OPERATOR_TOKEN,
	Message: "The lexeme did not contains a valid right-operator literal token, found [%s] token (%d:%d)",
}

// ErrorCoralSyntaxTreeInvalidRegexOperatorToken is returned when a lexeme
// does not contain a valid regex token in syntax tree.
var ErrorCoralSyntaxTreeInvalidRegexOperatorToken = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_INVALID_REGEX_OPERATOR_TOKEN,
	Message: "The lexeme did not contains a valid compilable RegEx token, found [%s] token",
}

// ErrorCoralSyntaxTreeMissingThenToken is returned when a lexeme does not
// contain a THEN token in syntax tree.
var ErrorCoralSyntaxTreeMissingThenToken = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_MISSING_THEN_TOKEN,
	Message: "The lexeme did not contains [THEN] token, found [%s] token (%d:%d)",
}

// ErrorCoralSyntaxTreeMissingReturnToken is returned when a lexeme does not
// contain an EXEC or RETURN token in syntax tree.
var ErrorCoralSyntaxTreeMissingReturnToken = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_MISSING_RETURN_TOKEN,
	Message: "The lexeme did not contains [EXEC] or [RETURN] token, found [%s] token (%d:%d)",
}

// ErrorCoralSyntaxTreeMissingExecLiteralToken is returned when a lexeme does
// not contain a literal token for EXEC branch in syntax tree.
var ErrorCoralSyntaxTreeMissingExecLiteralToken = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_MISSING_EXEC_LITERAL_TOKEN,
	Message: "The lexeme did not contains [LITERAL] token for EXEC branch, found [%s] token (%d:%d)",
}

// ErrorSyntaxTreeMissingExecLiteralToken is returned when a lexeme does not
// contain a literal token for EXEC branch in syntax tree.
var ErrorCoralSyntaxTreeMissingReturnLiteralToken = RemoraError{
	Code:    CODE_CORAL_SYNTAXTREE_MISSING_RETURN_LITERAL_TOKEN,
	Message: "The lexeme did not contains [LITERAL] token for RETURN branch, found [%s] token (%d:%d)",
}

// ErrorCoralOperatorInvalidBooleanConversion is returned when boolean
// conversion fails in operator evaluation.
var ErrorCoralOperatorInvalidBooleanConversion = RemoraError{
	Code:    CODE_CORAL_OPERATOR_INVALID_BOOLEAN_CONVERSION,
	Message: "Impossible to convert the value [%s] in boolean type",
}

// ErrorCoralOperatorInvalidStringConversion is returned when string
// conversion fails in operator evaluation.
var ErrorCoralOperatorInvalidStringConversion = RemoraError{
	Code:    CODE_CORAL_OPERATOR_INVALID_STRING_CONVERSION,
	Message: "Impossible to convert the value [%s] in string type",
}

// ErrorCoralOperatorInvalidNumericConversion is returned when numeric
// conversion fails in operator evaluation.
var ErrorCoralOperatorInvalidNumericConversion = RemoraError{
	Code:    CODE_CORAL_OPERATOR_INVALID_NUMERIC_CONVERSION,
	Message: "Impossible to convert the value [%s] in numeric type",
}

// ErrorCoralRunnerUnmarshalFailed is returned when unmarshal on
// compressed content fails.
var ErrorCoralRunnerUnmarshalFailed = RemoraError{
	Code:    CODE_CORAL_RUNNER_UNMARSHAL_FAILED,
	Message: "Impossible to correctly unmarshal content. Error: [%s]",
}

// ErrorCoralRunnerRawBodyResolutionFailed is returned when resolution
// fails on raw body content.
var ErrorCoralRunnerRawBodyResolutionFailed = RemoraError{
	Code:    CODE_CORAL_RUNNER_RAWBODY_RESOLUTION_FAILED,
	Message: "Impossible to apply resolution on raw body [%s]",
}

// ErrorCoralRunnerMappableBodyResolutionFailed is returned when resolution
// fails on mappable body content.
var ErrorCoralRunnerMappableBodyResolutionFailed = RemoraError{
	Code:    CODE_CORAL_RUNNER_MAPPABLEBODY_RESOLUTION_FAILED,
	Message: "Impossible to apply resolution on mappable body [%s]",
}

// ErrorCoralRunnerHeadersResolutionFailed is returned when resolution
// fails on headers content.
var ErrorCoralRunnerHeadersResolutionFailed = RemoraError{
	Code:    CODE_CORAL_RUNNER_HEADERS_RESOLUTION_FAILED,
	Message: "Impossible to apply resolution on header [%s]",
}
