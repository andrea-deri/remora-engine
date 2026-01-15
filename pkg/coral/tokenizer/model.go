package tokenizer

// TokenType represents the enumeration of valid tokens that the tokenizer can extract.
type TokenType string

// Token represents the structure of a new token extracted from statement.
type Token struct {
	Type  TokenType
	Value string
}

const (

	// TOKEN_IF define the type of 'if' static tokens
	TOKEN_IF TokenType = "IF"

	// TOKEN_THEN define the type of 'then' static tokens
	TOKEN_THEN TokenType = "THEN"

	// TOKEN_RETURN define the type of 'return' static tokens
	TOKEN_RETURN TokenType = "RETURN"

	// TOKEN_EXEC define the type of 'exec' static tokens
	TOKEN_EXEC TokenType = "EXEC"

	// TOKEN_OTHERWISE define the type of 'otherwise' static tokens
	TOKEN_OTHERWISE TokenType = "OTHERWISE"

	// TOKEN_DEFAULT define the type of 'default' static tokens. This is a
	// syntactic sugar for 'otherwise', used in default-only rules.
	TOKEN_DEFAULT TokenType = "DEFAULT"

	// TOKEN_END_BEHAVIOR define the type of behavior's end tokens
	TOKEN_END_BEHAVIOR TokenType = "END"

	// TOKEN_END_STATEMENT define the type of statement's end token
	TOKEN_END_STATEMENT TokenType = "END_STATEMENT"

	// TOKEN_LPAREN define the type of left parenthesys tokens
	TOKEN_LPAREN TokenType = "LPAREN"

	// TOKEN_RPAREN define the type of right parenthesys tokens
	TOKEN_RPAREN TokenType = "RPAREN"

	// TOKEN_UNARY_OPERATOR define the type of unary operation's tokens
	TOKEN_UNARY_OPERATOR TokenType = "UNARY_OPERATOR"

	// TOKEN_BINARY_OPERATOR define the type of binary operation's tokens
	TOKEN_BINARY_OPERATOR TokenType = "BINARY_OPERATOR"

	// TOKEN_LITERAL define the type of literal-value's tokens
	TOKEN_LITERAL TokenType = "LITERAL"

	// TOKEN_AND define the type of 'and' static tokens
	TOKEN_AND TokenType = "AND_CONJUNCTION"

	// TOKEN_OR define the type of 'or' static tokens
	TOKEN_OR TokenType = "OR_CONJUNCTION"
)

const (

	// RAW_TOKEN_IF represents the enumerative raw value of a IF token
	RAW_TOKEN_IF = "if"

	// RAW_TOKEN_THEN represents the enumerative raw value of a THEN token
	RAW_TOKEN_THEN = "then"

	// RAW_TOKEN_RETURN represents the enumerative raw value of a RETURN token
	RAW_TOKEN_RETURN = "return"

	// RAW_TOKEN_EXEC represents the enumerative raw value of a EXEC token
	RAW_TOKEN_EXEC = "exec"

	// RAW_TOKEN_OTHERWISE represents the enumerative raw value of a OTHERWISE token
	RAW_TOKEN_OTHERWISE = "otherwise"

	// RAW_TOKEN_DEFAULT represents the enumerative raw value of a DEFAULT token
	RAW_TOKEN_DEFAULT = "default"

	// RAW_TOKEN_END represents the enumerative raw value of a END token
	RAW_TOKEN_END = "end"

	// RAW_TOKEN_AND represents the enumerative raw value of a AND_CONJUNCTION token
	RAW_TOKEN_AND = "and"

	// RAW_TOKEN_OR represents the enumerative raw value of a OR_CONJUNCTION token
	RAW_TOKEN_OR = "or"

	// RAW_TOKEN_LITERAL_DELIMITER represents the enumerative raw value of the
	// delimiter for LITERAL tokens
	RAW_TOKEN_LITERAL_DELIMITER = '`'

	// RAW_TOKEN_LPAREN represents the enumerative raw value of a LPAREN token
	RAW_TOKEN_LPAREN = '('

	// RAW_TOKEN_RPAREN represents the enumerative raw value of a RPAREN token
	RAW_TOKEN_RPAREN = ')'
)
