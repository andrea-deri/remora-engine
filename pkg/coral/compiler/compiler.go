package compiler

import (
	"fmt"
	"remora/pkg/coral/parser"
	"remora/pkg/coral/profile"
	"remora/pkg/coral/tokenizer"
	"remora/pkg/protocol/transfer/btp"

	"github.com/vmihailenco/msgpack/v5"
)

// Compile converts a Minnow raw profile into an executable behavior Profile.
//
// For each behavior in the profile, the function serializes its effect, converts the behavior
// into a normalized rule expression using generateStatement and concatenates all generated
// statements into a complete behavior program.
//
// The final aggregated expression is then compiled into a runtime Profile object by the internal
// behavior compiler. If any serialization step fails, the function panics as such a condition indicates
// a non-recoverable configuration error.
func Compile(minnow btp.MinnowProfile) (profile.Profile, error) {

	var completeProfileStatement string
	for _, behavior := range minnow.Behaviors {

		behaviorStatement := generateBehaviorStatement(behavior)
		completeProfileStatement = fmt.Sprintf("%s%s", completeProfileStatement, behaviorStatement)
	}

	return compileBehavior(completeProfileStatement)
}

// generateBehaviorStatement converts a single Minnow behavior into its textual rule representation.
//
// The function serializes the behavior's Effect using MessagePack, determines the appropriate
// operation token based on the effect type (return, action, default) and embeds the serialized
// payload between literal delimiters.
//
// The resulting rule has the form:
//
//	<condition> <operation> <delimiter><binary-payload><delimiter> END
//
// This normalized representation is suitable for concatenation into the full behavior program
// consumed by the rule compiler.
func generateBehaviorStatement(behavior btp.MinnowBehavior) string {

	compressedBehavior, err := msgpack.Marshal(behavior.Effect)
	if err != nil {
		panic(err)
	}

	var operation string
	switch behavior.Effect.Type {
	case "return":
		operation = fmt.Sprintf("%s %s", tokenizer.RAW_TOKEN_THEN, tokenizer.RAW_TOKEN_RETURN)
	case "action":
		operation = fmt.Sprintf("%s %s", tokenizer.RAW_TOKEN_THEN, tokenizer.RAW_TOKEN_EXEC)
	case "default":
		operation = tokenizer.RAW_TOKEN_RETURN
	}

	literalDelimiter := string(tokenizer.RAW_TOKEN_LITERAL_DELIMITER)
	return fmt.Sprintf("%s %s %s%s%s %s ", behavior.Condition, operation, literalDelimiter, compressedBehavior, literalDelimiter, tokenizer.RAW_TOKEN_END)
}

// compileBehavior converts a raw behavior statement into an executable Profile.
//
// The function performs two sequential steps:
//  1. Tokenization: the input behavior statement is transformed into a stream of tokens
//     compatible with the behavior grammar.
//  2. Parsing: the token stream is parsed into a Profile, which represents the compiled
//     behavioral program ready to be executed.
//
// The returned Profile can later be evaluated against an input map, enabling full Minnow
// behavior validation and execution. If tokenization or parsing fails, an error is returned
// and no Profile is produced.
func compileBehavior(behaviorStatement string) (profile.Profile, error) {

	tokenizedStatement, err := tokenizer.TokenizeStatement(behaviorStatement)
	if err != nil {
		return profile.Profile{}, err
	}

	behaviorProfile, err := parser.Parse(tokenizedStatement)
	return behaviorProfile, err
}
