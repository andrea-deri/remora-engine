package runner

import (
	"remora/pkg/coral/profile"
	"remora/pkg/customerror"
	"remora/pkg/protocol/transfer/btp"

	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

// Evaluate applies a behavior profile to an input map and determines the first behavior
// whose condition is satisfied.
//
// Behaviors are evaluated sequentially. As soon as a condition evaluates to true, its associated
// effect is executed and evaluation stops. If the effect defines a RETURN clause, its value is
// returned, otherwise an ACTION clause is executed. If no behavior is satisfied, tipically a
// DEFAULT behavior (i.e. the last defined in the profile) is always applied and can be used as
// parachute, in order to avoid strange behaviors.
//
// The function returns:
//   - the value produced by the applied effect (if any),
//   - a boolean indicating whether at least one behavior was satisfied,
//   - an error if condition evaluation or effect execution fails.
func Evaluate(profile profile.Profile, input map[string]any) (any, bool, error) {

	for _, behavior := range profile.Behaviors {

		isInputCompliantToBehavior, err := isConditionSatisfied(behavior.Condition, input)
		if err != nil {
			return nil, false, err
		}

		if isInputCompliantToBehavior {

			if behavior.Effect.Return != nil {
				valueFromEffect, err := applyReturn(behavior, input)
				log.Trace().Msgf("Applied return effect: [%v]", valueFromEffect)
				if err != nil {
					return nil, false, err
				}
				return valueFromEffect, true, nil
			}

			if behavior.Effect.Action != nil {
				returnValue, err := applyAction(behavior, input)
				log.Trace().Msgf("Applied action effect: [%v]", returnValue)
				if err != nil {
					return nil, false, err
				}
				return returnValue, true, nil
			}
		}
	}

	// Hypotetically, this is never reached if a default behavior is set, but it
	// is defined as fallback to avoid strange situations.
	return nil, false, nil
}

// applyReturn evaluates and materializes the RETURN effect of a compliant behavior.
//
// The function decodes the serialized effect definition, resolves its content by injecting
// values from the input map and produces a concrete EffectResult.
// Resolution supports:
//   - constant values, not injecting and returning the as-is response
//   - static values, by injecting values from the input object map into the response
//   - dynamic values, by injecting values calculated at runtime or retrievable from the local state
//
// An error is returned if the effect cannot be decoded or if value resolution fails.
func applyReturn(behavior profile.Behavior, input map[string]any) (profile.EffectResult, error) {

	behaviorReturnEffect := behavior.Effect.Return.Value

	var effect btp.MinnowBehaviorEffect
	err := msgpack.Unmarshal([]byte(behaviorReturnEffect), &effect)
	if err != nil {
		return profile.EffectResult{}, customerror.NewError(customerror.ErrorCoralRunnerUnmarshalFailed, err)
	}
	log.Trace().Msgf("Found effect from behavior: [%v]", effect)

	headers, err := applyResolutionOnMappableBody(effect.Headers, input)
	if err != nil {
		return profile.EffectResult{}, customerror.NewError(customerror.ErrorCoralRunnerRawBodyResolutionFailed, effect.Headers)
	}

	response := profile.EffectResult{
		StatusCode:  effect.StatusCode,
		ContentType: effect.ContentType,
		Headers:     headers,
	}

	if effect.RawBody != "" {

		plainStringBody, err := applyResolutionOnRawBody(effect.RawBody, input)
		if err != nil {
			return profile.EffectResult{}, customerror.NewError(customerror.ErrorCoralRunnerRawBodyResolutionFailed, effect.RawBody)
		}
		response.RawBody = plainStringBody

	} else {

		mappableBody, err := applyResolutionOnMappableBody(effect.MappableBody, input)
		if err != nil {
			return profile.EffectResult{}, customerror.NewError(customerror.ErrorCoralRunnerRawBodyResolutionFailed, effect.MappableBody)
		}
		response.MappableBody = mappableBody
	}

	return response, nil
}

// TODO: currently not implemented
func applyAction(behavior profile.Behavior, input map[string]any) (any, error) {

	/*
	   TODO
	   - perform operation defined in Action
	     - open a sandbox
	     - execute code defined in action
	     - generate an output map containing dynamically generated data
	   - deconvert Body from GZIP
	   - perform resolution if plain text body contains at least 1 occurrence of '${'
	     - if it contains "resolvable fields”, extract all fields
	     - for each field:
	         - if it is in the form ‘${actout.value}’, injection occurs from the action output map
	         - if it is in the form ‘${some.value}’, injection occurs from the input map
	   	  - if it is in the form ‘${!someFunct}’, the injection occurs by executing a native function
	   - return object in the form of a map with:
	     - status code
	     - content type
	     - various headers
	     - body with resolution completed
	*/
	return nil, customerror.NewError(customerror.ErrorCoralSyntaxTreeFunctionalityNotImplemented, "action")
}
