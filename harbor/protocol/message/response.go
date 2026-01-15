package message

import (
	"remora/pkg/customerror"
)

// BodyType defines the type of body in the HarborResponse. This is an alias for string to
// improve code readability and set predefined constants.
// The accepted values are: BODY_TYPE_RAW, BODY_TYPE_MAP.
type BodyType int

const (

	// BODY_TYPE_RAW defines raw string body type
	BODY_TYPE_RAW BodyType = 0

	// BODY_TYPE_MAP defines mappable object body type
	BODY_TYPE_MAP BodyType = 1
)

// HarborResponse represents the response message sent from Dispatcher to Proxy.
// It contains all the necessary information about the HTTP response to be sent back to the client.
type HarborResponse struct {

	// ContentType indicates the type of content in the response. The value is defined
	// using the content type passed in the request or the one set by user in the
	// behavior's content the one as fallback value.
	ContentType string

	// Body contains the body of the response. It can be either a raw string
	// or a mappable object, depending on the BodyType field.
	Body any

	// Headers contains the HTTP headers for the response. The value is defined
	// using both the headers set by user in the behavior's content and the
	// standard headers set by REMORA engine.
	Headers map[string]string

	// StatusCode is the HTTP status code of the response. The value is defined
	// using the status code set by user in the behavior's content.
	StatusCode int

	// BodyType indicates the type of body in the response. It can be either
	// BODY_TYPE_RAW (for raw string body) or BODY_TYPE_MAP for (mappable object body).
	BodyType BodyType
}

// BuildErrorResponse constructs consistent error responses across the application with
// standardized error response body. It takes an error code and a message as input and
// returns a map[string]any containing the error details.
func BuildErrorResponse(errorCode customerror.RemoraErrorCode, message string) map[string]any {

	return map[string]any{
		"code":    errorCode,
		"message": message,
	}
}
