package parser

import (
	"remora/pkg/coral/tokenizer"
	"slices"
)

// TokenType defines a lexeme, i.e. the type of a token in the syntax tree that could be related to a single rule.
type Lexeme struct {

	// Tokens contains all tokens belonging to the lexeme.
	Tokens []tokenizer.Token

	// cursor is the position of the handled token in the analyzed lexeme.
	cursor int

	// statementCursor is the index of the lexeme in the entire parser.
	statementCursor int
}

// token returns the current token in the lexeme based on the cursor position. If the cursor exceeds
// the available token range, the function returns an END token, ensuring callers always receive a valid
// token object without needing explicit bounds checks.
func (lexeme *Lexeme) token() tokenizer.Token {

	if lexeme.cursor > len(lexeme.Tokens) {
		return tokenizer.Token{Type: tokenizer.TOKEN_END_BEHAVIOR}
	}
	return lexeme.Tokens[lexeme.cursor]
}

// typedToken retrieves the current token and checks whether its type matches one of the expected
// token types. This helper abstracts type-checking logic used heavily during parsing, allowing
// callers to validate grammar structure without duplicating comparisons.
//
// It returns the retrieved token and a boolean indicating whether the token type is compliant
// with any of the expected types.
func (lexeme *Lexeme) typedToken(expectedTypes ...tokenizer.TokenType) (tokenizer.Token, bool) {

	currentToken := lexeme.token()
	if !slices.Contains(expectedTypes, currentToken.Type) {
		return currentToken, false
	}
	return currentToken, true
}

// next advances the internal cursor to the next token in the lexeme. If the cursor moves past
// the end of the token list, it resets the cursor to zero. This behavior provides a fallback
// in order to avoid cursor overflow in parsing routines.
func (lexeme *Lexeme) next() {

	if lexeme.cursor > len(lexeme.Tokens) {
		lexeme.cursor = 0
	} else {
		lexeme.cursor++
	}
}
