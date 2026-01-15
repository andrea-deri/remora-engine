package parser

import (
	"regexp"
	"slices"
	"strings"

	"remora/pkg/coral/profile"
	"remora/pkg/coral/spec/operator"
	"remora/pkg/coral/tokenizer"
	"remora/pkg/customerror"
)

// Parse transforms a sequence of raw tokens into an executable Profile. It coordinates the entire
// parsing process: it divides the sequence into lexemes, validates each block, constructs Behaviors
// through syntactic analysis and returns a Profile structure ready for execution by the Engine.
//
// It returns an error if the token structure is malformed or if a Behavior cannot be generated correctly.
func Parse(statementTokens []tokenizer.Token) (profile.Profile, error) {

	lexemes, err := extractLexemes(statementTokens)
	if err != nil {
		return profile.Profile{}, err
	}

	var behaviors []profile.Behavior
	for _, lexeme := range lexemes {
		behavior, err := parseBehavior(&lexeme)
		if err != nil {
			return profile.Profile{}, err
		}
		behaviors = append(behaviors, behavior)
	}

	return profile.Profile{Behaviors: behaviors}, nil
}

// extractLexemes groups tokens into logical units (the lexemes, representing a single DSL behavior).
// The function calls generateLexemes to create the blocks and then validateLexeme to ensure their
// syntactic correctness.
//
// It returns an error if a block does not comply with the required structure (invalid start or end token).
func extractLexemes(statementToken []tokenizer.Token) ([]Lexeme, error) {

	lexemes := generateLexemes(statementToken)

	for index, lexeme := range lexemes {
		err := validateLexeme(lexeme, index)
		if err != nil {
			return nil, err
		}
	}
	return lexemes, nil
}

// validateLexeme checks that a single lexeme is structurally valid. It checks that it starts with
// an allowed token (IF/OTHERWISE/DEFAULT) and ends with an END token.
//
// The function does not interpret the content, but ensures that the block boundaries comply with the high-level grammar.
//
// It returns an error if the lexeme form is malformed.
func validateLexeme(lexeme Lexeme, lexemeIndex int) error {

	startToken := lexeme.Tokens[0]
	endToken := lexeme.Tokens[len(lexeme.Tokens)-1]

	isStartTokenValid := startToken.Type != tokenizer.TOKEN_IF && startToken.Type != tokenizer.TOKEN_OTHERWISE && startToken.Type != tokenizer.TOKEN_DEFAULT
	isEndTokenValid := endToken.Type != tokenizer.TOKEN_END_BEHAVIOR

	if isStartTokenValid || isEndTokenValid {
		return customerror.NewError(customerror.ErrorCoralSyntaxTreeMalformedLexeme, lexemeIndex, startToken.Type, endToken.Type)
	}
	return nil
}

// generateLexemes performs primary segmentation of tokens into lexemes. Each occurrence of
// IF/OTHERWISE/DEFAULT tokens opens a new lexeme and END token concludes the current block.
//
// The function does not perform validations and does not interpret the internal content.
// It only produces a sequential partition of tokens, preserving the original order.
func generateLexemes(statementToken []tokenizer.Token) []Lexeme {

	currentLexemeCursor := 0
	lexemes := []Lexeme{}

	for _, token := range statementToken {

		if token.Type == tokenizer.TOKEN_IF || token.Type == tokenizer.TOKEN_OTHERWISE || token.Type == tokenizer.TOKEN_DEFAULT {
			lexemes = append(lexemes, Lexeme{statementCursor: currentLexemeCursor})
		}

		currentLexeme := &lexemes[currentLexemeCursor]
		currentLexeme.Tokens = append(currentLexeme.Tokens, token)

		if token.Type == tokenizer.TOKEN_END_BEHAVIOR {
			currentLexemeCursor++
		}
	}
	return lexemes
}

// parseBehavior constructs a complete Behavior from the given lexeme. It interprets the structure
// of a single block: it reads the initial token, extracts the condition after IF token (or find
// fallback condition after DEFAULT/OTHERWISE token), interprets the THEN section, then checks for
// the presence of the concluding END token. It is responsible for the semantic correctness of the
// block and returns a Behavior ready for execution.
//
// It returns an error if the token sequence does not follow the required grammar.
func parseBehavior(lexeme *Lexeme) (profile.Behavior, error) {

	// Check if start token is 'IF', 'OTHERWISE' or 'DEFAULT' token
	currentToken, isValid := lexeme.typedToken(tokenizer.TOKEN_IF, tokenizer.TOKEN_OTHERWISE, tokenizer.TOKEN_DEFAULT)
	if !isValid {
		return profile.Behavior{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeInvalidStartingTokenLexeme, currentToken.Type, lexeme.statementCursor, lexeme.cursor)
	}
	lexeme.next()

	// Set condition with "otherwise" case as fallback if current token is not 'IF'
	condition := profile.Condition{
		Type:   operator.EXPR_TYPE_CLAUSE,
		Clause: &profile.Clause{Operator: operator.SPECIAL_DEFAULT},
	}

	isIFToken := currentToken.Type == tokenizer.TOKEN_IF
	if isIFToken {
		orCondition, err := parseConditionFromDisjunction(lexeme)
		if err != nil {
			return profile.Behavior{}, err
		}
		condition = orCondition
	}

	effect, err := parseEffect(lexeme, isIFToken)
	if err != nil {
		return profile.Behavior{}, err
	}
	lexeme.next()

	// Check if end token is 'END' token
	currentToken, isValid = lexeme.typedToken(tokenizer.TOKEN_END_BEHAVIOR)
	if !isValid {
		return profile.Behavior{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeInvalidEndingTokenLexeme, currentToken.Type, lexeme.statementCursor, lexeme.cursor)
	}

	behavior := profile.Behavior{
		Condition: condition,
		Effect:    effect,
	}
	return behavior, nil
}

// parseConditionFromDisjunction interprets a logical condition constructed as an OR expression.
// It analyzes a first AND term, then detects any additional contiguous OR terms.
//
// If multiple terms are present, it creates a logical OR node, otherwise it directly returns the first term.
//
// It returns an error if conditions is malformed or there are unexpected operators.
func parseConditionFromDisjunction(lexeme *Lexeme) (profile.Condition, error) {

	firstCondition, err := parseConditionFromConjunction(lexeme)
	if err != nil {
		return profile.Condition{}, err
	}
	conditions := []profile.Condition{firstCondition}

	// Every time an OR token is found, add the OR-disjuncted term to the slice
	for {

		_, isORToken := lexeme.typedToken(tokenizer.TOKEN_OR)
		if !isORToken {
			break
		}
		lexeme.next()

		conditionFromLowerLevel, err := parseConditionFromConjunction(lexeme)
		if err != nil {
			return profile.Condition{}, err
		}
		conditions = append(conditions, conditionFromLowerLevel)
	}

	// The only elaborated term is the first, no other element were found in OR-disjuncted terms
	if len(conditions) == 1 {
		return firstCondition, nil
	}

	condition := profile.Condition{
		Type:  operator.EXPR_TYPE_OR,
		Terms: conditions,
	}
	return condition, nil
}

// parseConditionFromConjunction processes a logical condition based on AND expression. It reads a
// first term (clause or sub-expression), then adds all subsequent terms connected by AND.
//
// If multiple terms are present, it create a single logical AND node, otherwise it directly returns the base condition.
//
// It returns an error if the structure is invalid.
func parseConditionFromConjunction(lexeme *Lexeme) (profile.Condition, error) {

	// Extract the first clause-conjuncted term
	firstCondition, err := parseCondition(lexeme)
	if err != nil {
		return profile.Condition{}, err
	}
	conditions := []profile.Condition{firstCondition}

	// Every time an AND token is found, add the AND-conjuncted term to the slice
	for {

		_, isANDToken := lexeme.typedToken(tokenizer.TOKEN_AND)
		if !isANDToken {
			break
		}
		lexeme.next()

		conditionFromSameLevel, err := parseCondition(lexeme)
		if err != nil {
			return profile.Condition{}, err
		}
		conditions = append(conditions, conditionFromSameLevel)
	}

	// The only elaborated term is the first, no other element were found in AND-conjuncted terms
	if len(conditions) == 1 {
		return firstCondition, nil
	}

	condition := profile.Condition{
		Type:  operator.EXPR_TYPE_AND,
		Terms: conditions,
	}
	return condition, nil
}

// parseCondition interprets a single condition element. This can be a parenthetical expression,
// which is parsed recursively as an OR expression, or an atomic clause generated by parseClause.
//
// It handles the correct closing of parentheses and reports syntax errors in case of mismatches.
func parseCondition(lexeme *Lexeme) (profile.Condition, error) {

	_, isLeftParenthesys := lexeme.typedToken(tokenizer.TOKEN_LPAREN)
	if isLeftParenthesys {

		lexeme.next()

		innerORCondition, err := parseConditionFromDisjunction(lexeme)
		if err != nil {
			return profile.Condition{}, err
		}

		_, isRightParenthesys := lexeme.typedToken(tokenizer.TOKEN_RPAREN)
		if !isRightParenthesys {
			return profile.Condition{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeMissingRightParenthesys, lexeme.statementCursor, lexeme.cursor)
		}
		lexeme.next()

		return innerORCondition, nil
	}

	clause, err := parseClause(lexeme)
	if err != nil {
		return profile.Condition{}, err
	}

	condition := profile.Condition{
		Type:   operator.EXPR_TYPE_CLAUSE,
		Clause: &clause,
	}
	return condition, nil
}

// parseClause constructs an atomic clause of the DSL. It extracts a field (literal), an operator (unary or binary)
// and a right-hand side value that can be a simple literal or a compiled regex, depending on the operator.
//
// The function applies the semantic rules of the Coral DSL and returns an error if any of the tokens aren't
// compliant with the expected form or if a regex is invalid.
func parseClause(lexeme *Lexeme) (profile.Clause, error) {

	fieldToken, isValidLiteral := lexeme.typedToken(tokenizer.TOKEN_LITERAL)
	if !isValidLiteral {
		return profile.Clause{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeMissingLiteralToken, fieldToken.Type, lexeme.statementCursor, lexeme.cursor)
	}
	fieldToken.Value = strings.ToLower(fieldToken.Value)
	lexeme.next()

	operatorToken, isValidOperator := lexeme.typedToken(tokenizer.TOKEN_UNARY_OPERATOR, tokenizer.TOKEN_BINARY_OPERATOR)
	if !isValidOperator {
		return profile.Clause{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeInvalidOperatorToken, fieldToken.Type, lexeme.statementCursor, lexeme.cursor)
	}
	operatorContent := operator.Operator(operatorToken.Value)
	lexeme.next()

	var valueContent any
	if operator.BINARY_REGEX == operatorContent {

		valueToken, isValidLiteral := lexeme.typedToken(tokenizer.TOKEN_LITERAL)
		if !isValidLiteral {
			return profile.Clause{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeInvalidRightOperatorToken, fieldToken.Type, lexeme.statementCursor, lexeme.cursor)
		}

		compiledRegex, err := regexp.Compile(valueToken.Value)
		if err != nil {
			return profile.Clause{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeInvalidRegexOperatorToken, valueToken.Value)
		}
		valueContent = compiledRegex

		lexeme.next()

	} else if slices.Contains(operator.BinaryOperators, operatorContent) {

		valueToken, isValidLiteral := lexeme.typedToken(tokenizer.TOKEN_LITERAL)
		if !isValidLiteral {
			return profile.Clause{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeInvalidRightOperatorToken, fieldToken.Type, lexeme.statementCursor, lexeme.cursor)
		}
		valueContent = valueToken.Value

		lexeme.next()
	}

	clause := profile.Clause{
		Field:    fieldToken.Value,
		Operator: operatorContent,
		Value:    valueContent,
	}
	return clause, nil
}

// parseEffect constructs the effect section (THEN-EXEC/RETURN). If the rule is introduced by IF,
// it checks for the existence of the THEN token. It then interprets the type of effect (EXEC or RETURN)
// and reads its literal argument.
//
// It returns a fully structured Effect object or an error if the token sequence is not compliant with the grammar.
func parseEffect(lexeme *Lexeme, isIFToken bool) (profile.Effect, error) {

	// Check compliance to [IF|OTHERWISE|DEFAULT]-THEN structure
	if isIFToken {
		currentToken, isTHENToken := lexeme.typedToken(tokenizer.TOKEN_THEN)
		if !isTHENToken {
			return profile.Effect{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeMissingThenToken, currentToken.Type, lexeme.statementCursor, lexeme.cursor)
		}
		lexeme.next()
	}

	currentToken, isValid := lexeme.typedToken(tokenizer.TOKEN_EXEC, tokenizer.TOKEN_RETURN)
	if !isValid {
		return profile.Effect{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeMissingReturnToken, currentToken.Type, lexeme.statementCursor, lexeme.cursor)
	}
	lexeme.next()

	switch currentToken.Type {

	case tokenizer.TOKEN_EXEC:
		currentToken, isValidLiteral := lexeme.typedToken(tokenizer.TOKEN_LITERAL)
		if !isValidLiteral {
			return profile.Effect{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeMissingExecLiteralToken, currentToken.Type, lexeme.statementCursor, lexeme.cursor)
		}
		effect := profile.Effect{
			Action: &profile.Action{
				Code: currentToken.Value,
			},
		}
		return effect, nil

	case tokenizer.TOKEN_RETURN:
		currentToken, isValidLiteral := lexeme.typedToken(tokenizer.TOKEN_LITERAL)
		if !isValidLiteral {
			return profile.Effect{}, customerror.NewError(customerror.ErrorCoralSyntaxTreeMissingReturnLiteralToken, currentToken.Type, lexeme.statementCursor, lexeme.cursor)
		}
		effect := profile.Effect{
			Return: &profile.Return{
				Value: currentToken.Value,
			},
		}
		return effect, nil
	}

	return profile.Effect{}, nil
}
