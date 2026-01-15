package runner

import (
	"fmt"
	"remora/pkg/conversion/format"
	"remora/pkg/conversion/gzip"
	"remora/pkg/conversion/maps"
	"remora/pkg/conversion/resolution"
	"strings"

	"github.com/rs/zerolog/log"
)

// applyResolutionOnRawBody resolves all resolution placeholders contained in a raw string body
// using values derived from the provided input map.
//
// A resolution placeholder is any expression wrapped in the "${...}" syntax. Placeholders may
// reference input fields directly or define function expressions using the
// "!functionName(parameter1, parameter2, ...)" format, which are evaluated at runtime.
//
// The function transparently handles GZIP-compressed bodies and always operates on the
// decompressed representation. If no resolution placeholders are present, the body is
// returned unchanged.
func applyResolutionOnRawBody(rawBody string, input map[string]any) (string, error) {

	updatedRawBody, err := gzip.CompressedStringToPlainString(rawBody)
	if err != nil {
		log.Error().Msgf("An error occurred while decompressing string from GZIP: [%v]", err)
		return "", err
	}

	if strings.Contains(updatedRawBody, "${") {

		resolutionPlaceholderFields := resolution.GetResolvableFieldsFromString(updatedRawBody)
		log.Trace().Msgf("Found [%d] values to be resolved!", len(resolutionPlaceholderFields))

		for _, fieldWithResolutionPlaceholder := range resolutionPlaceholderFields {

			var resolvedValue any
			if strings.HasPrefix(fieldWithResolutionPlaceholder, "!") {

				resolvedValue = applyResolutionOnFunction(input, fieldWithResolutionPlaceholder)
				if resolvedValue != nil {
					contentToBeResolved := fmt.Sprintf("${%v}", fieldWithResolutionPlaceholder)
					contentAsReplacement := fmt.Sprintf("%v", resolvedValue)
					updatedRawBody = strings.ReplaceAll(updatedRawBody, contentToBeResolved, contentAsReplacement)
				}

			} else {

				resolvedValue = maps.ExtractFieldValue(input, fieldWithResolutionPlaceholder)
				if resolvedValue != nil {
					contentToBeResolved := fmt.Sprintf("${%s}", fieldWithResolutionPlaceholder)
					contentAsReplacement := resolvedValue.(string)
					updatedRawBody = strings.ReplaceAll(updatedRawBody, contentToBeResolved, contentAsReplacement)
				}
			}

			log.Trace().Msgf("Resolvable field [%v] in return static body is resolved with value [%v].", fieldWithResolutionPlaceholder, resolvedValue)
		}
	}
	return updatedRawBody, nil
}

// applyResolutionOnMappableBody resolves all resolution placeholders contained in a mappable body
// using values derived from the provided input map.
//
// A resolution placeholder is any value encoded using the "${...}" syntax. Placeholders may
// reference input fields or define function expressions using the "!functionName(...)" form.
// Function expressions are evaluated and resolved recursively.
//
// The function operates on a deep copy of the original map to avoid mutating the input data.
// If the provided mappable body is empty, an empty map is returned.
func applyResolutionOnMappableBody(mappableBody map[string]any, input map[string]any) (map[string]any, error) {

	if len(mappableBody) == 0 {
		return make(map[string]any), nil
	}

	// Execute a deep copy in order to avoid applying changes on Minnow behavior effect's body
	updatedMappableBody := maps.DeepCopy(mappableBody)

	resolutionPlaceholderFields := resolution.GetResolvableFieldNamesFromMap(updatedMappableBody)
	for _, fieldWithResolutionPlaceholder := range resolutionPlaceholderFields {

		rawResolutionPlaceholder := maps.ExtractFieldValue(updatedMappableBody, fieldWithResolutionPlaceholder)

		resolutionPlaceholder := ""
		if rawResolutionPlaceholder != nil {
			resolutionPlaceholder = rawResolutionPlaceholder.(string)
		}
		resolvableFields := resolution.ExtractFieldsFromResolutionPlaceholders(resolutionPlaceholder)

		for _, resolvableField := range resolvableFields {

			var resolvedValue any
			if strings.HasPrefix(resolvableField, "!") {

				resolvedValue = applyResolutionOnFunction(input, resolvableField)
				if resolvedValue != nil {
					resolution.InsertResolvedValueInMappableBody(updatedMappableBody, resolvableField, fieldWithResolutionPlaceholder, resolvedValue)
				}

			} else {
				resolvedValue = maps.ExtractFieldValue(input, resolvableField)
				if resolvedValue != nil {
					resolution.InsertResolvedValueInMappableBody(updatedMappableBody, resolvableField, fieldWithResolutionPlaceholder, resolvedValue)
				}
			}

			log.Trace().Msgf("From resolvable field [%v] in object map, the key [%v] is resolved with value [%v].", fieldWithResolutionPlaceholder, resolutionPlaceholder, resolvedValue)
		}
	}
	return updatedMappableBody, nil
}

// applyResolutionOnFunction resolves and evaluates a function expression encoded in a resolvable field.
//
// The expression may reference input fields, constant values or nested function calls.
// All parameters are resolved first, then the function is executed and its result returned.
func applyResolutionOnFunction(input map[string]any, resolvableField string) any {

	functionName, functionParameters := resolution.GetResolvableFieldNamesFromFunction(resolvableField)
	log.Trace().Msgf("Found function [%v] with parameters %v that require a resolution process", functionName, functionParameters)

	resolvedParameters := make([]any, len(functionParameters))
	for i, functionParameter := range functionParameters {

		if strings.HasPrefix(functionParameter, "$") {

			parameter := format.RemoveFirstChar(functionParameter)
			resolvedParameters[i] = maps.ExtractFieldValue(input, parameter)

		} else if strings.HasPrefix(functionParameter, "!") {

			resolvedParameters[i] = applyResolutionOnFunction(input, functionParameter)
			log.Trace().Msgf("Found nested function [%s] that generated parameter [%v]", functionParameter, resolvedParameters[i])

		} else {

			if strings.HasPrefix(functionParameter, `"`) && strings.HasSuffix(functionParameter, `"`) {
				functionParameter = format.Unescape(functionParameter)
			}
			resolvedParameters[i] = functionParameter
		}
	}

	return resolution.FromFunction(functionName, resolvedParameters)
}
