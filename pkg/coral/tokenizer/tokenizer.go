package tokenizer

import (
	"fmt"
	"remora/pkg/coral/spec/operator"
	"slices"
	"strings"
	"unicode"
)

// TokenizeStatement converts a rule statement into a sequence of tokens.
// It scans the input string from left to right recognizing delimiters, parentheses,
// quoted literals, keywords, and operators, and returns the resulting token list.
// Whitespace is ignored.
//
// Any failure in literal extraction produces an error. All other malformed inputs
// produce erro inside helper routines.
func TokenizeStatement(statement string) ([]Token, error) {

	var tokens []Token
	statementSize := len(statement)

	cursor := 0
	for cursor < statementSize {

		character := statement[cursor]

		if unicode.IsSpace(rune(character)) {
			cursor++
			continue
		}

		switch character {

		case RAW_TOKEN_LPAREN:
			tokens = append(tokens, Token{Type: TOKEN_LPAREN})
			cursor++
			continue
		case RAW_TOKEN_RPAREN:
			tokens = append(tokens, Token{Type: TOKEN_RPAREN})
			cursor++
			continue
		case RAW_TOKEN_LITERAL_DELIMITER:
			token, updatedIndex, err := extractLiteralToken(statement, cursor)
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, token)
			cursor = updatedIndex
			continue
		}

		token, updatedIndex := extractToken(statement, cursor)

		cursor = updatedIndex
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// extractToken identifies the next non-literal token starting at tokenIndex.
// It reads characters until the next whitespace, converting the substring to a keyword,
// operator or literal token. Unknown words are treated as literal tokens. Empty words are ignored.
func extractToken(statement string, tokenIndex int) (Token, int) {

	statementSize := len(statement)
	wordEndIndex := tokenIndex

	// Extract required word until a space character is found
	for wordEndIndex < statementSize && !unicode.IsSpace(rune(statement[wordEndIndex])) {
		wordEndIndex++
	}

	lowercaseWord := strings.ToLower(statement[tokenIndex:wordEndIndex])

	var token Token
	switch {

	case lowercaseWord == RAW_TOKEN_IF:
		token = Token{Type: TOKEN_IF}
	case lowercaseWord == RAW_TOKEN_THEN:
		token = Token{Type: TOKEN_THEN}
	case lowercaseWord == RAW_TOKEN_EXEC:
		token = Token{Type: TOKEN_EXEC}
	case lowercaseWord == RAW_TOKEN_RETURN:
		token = Token{Type: TOKEN_RETURN}
	case lowercaseWord == RAW_TOKEN_OTHERWISE:
		token = Token{Type: TOKEN_OTHERWISE}
	case lowercaseWord == RAW_TOKEN_DEFAULT:
		token = Token{Type: TOKEN_DEFAULT}
	case lowercaseWord == RAW_TOKEN_END:
		token = Token{Type: TOKEN_END_BEHAVIOR}
	case lowercaseWord == RAW_TOKEN_AND:
		token = Token{Type: TOKEN_AND}
	case lowercaseWord == RAW_TOKEN_OR:
		token = Token{Type: TOKEN_OR}
	case slices.Contains(operator.UnaryOperators, operator.Operator(lowercaseWord)):
		token = Token{Type: TOKEN_UNARY_OPERATOR, Value: lowercaseWord}
	case slices.Contains(operator.BinaryOperators, operator.Operator(lowercaseWord)):
		token = Token{Type: TOKEN_BINARY_OPERATOR, Value: lowercaseWord}
	case lowercaseWord == "":
		// no-op
	default:
		token = Token{Type: TOKEN_LITERAL, Value: lowercaseWord}
	}

	return token, wordEndIndex + 1
}

// extractLiteralToken parses a quoted literal (delimited by the literal delimiter token) starting at
// tokenIndex. It scans forward to find the matching closing delimiter.
//
// If found, it returns the literal token and the updated cursor position. If no closing delimiter exists,
// an error is returned.
func extractLiteralToken(statement string, tokenIndex int) (Token, int, error) {

	statementSize := len(statement)

	// Search the index of closure character
	literalIndex := tokenIndex + 1
	for literalIndex < statementSize && statement[literalIndex] != RAW_TOKEN_LITERAL_DELIMITER {
		literalIndex++
	}

	if literalIndex >= statementSize {
		return Token{}, tokenIndex, fmt.Errorf("literal value not correctly escaped at index [%d] in statement", literalIndex)
	}

	literalToken := Token{
		Type:  TOKEN_LITERAL,
		Value: statement[tokenIndex+1 : literalIndex],
	}
	return literalToken, literalIndex + 1, nil
}
