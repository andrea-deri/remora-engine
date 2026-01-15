package message

import (
	"context"
)

// ContentType represents the type of content in the request or response.
// This is an alias for string to improve code readability and set predefined constants.
// The accepted values are: json, xml, string.
type ContentType string

const (

	// CONTENT_TYPE_JSON defines JSON content type
	CONTENT_TYPE_JSON ContentType = "json"

	// CONTENT_TYPE_XML defines XML content type
	CONTENT_TYPE_XML ContentType = "xml"

	// CONTENT_TYPE_STRING defines plain string content type
	CONTENT_TYPE_STRING ContentType = "string"
)

// HTTPMethod represents the HTTP method of the request. This is an alias for string to
// improve code readability and set predefined constants.
// The accepted values are: GET, POST, PUT, DELETE, OPTIONS, HEAD, PATCH, TRACE, CONNECT.
type HTTPMethod string

const (

	// HTTP_METHOD_GET defines GET HTTP method
	HTTP_METHOD_GET HTTPMethod = "get"

	// HTTP_METHOD_POST defines POST HTTP method
	HTTP_METHOD_POST HTTPMethod = "post"

	// HTTP_METHOD_PUT defines PUT HTTP method
	HTTP_METHOD_PUT HTTPMethod = "put"

	// HTTP_METHOD_DELETE defines DELETE HTTP method
	HTTP_METHOD_DELETE HTTPMethod = "delete"

	// HTTP_METHOD_OPTIONS defines OPTIONS HTTP method
	HTTP_METHOD_OPTIONS HTTPMethod = "options"

	// HTTP_METHOD_HEAD defines HEAD HTTP method
	HTTP_METHOD_HEAD HTTPMethod = "head"

	// HTTP_METHOD_PATCH defines PATCH HTTP method
	HTTP_METHOD_PATCH HTTPMethod = "patch"

	// HTTP_METHOD_PATCH defines PATCH HTTP method
	HTTP_METHOD_TRACE HTTPMethod = "trace"

	// HTTP_METHOD_CONNECT defines CONNECT HTTP method
	HTTP_METHOD_CONNECT HTTPMethod = "connect"
)

// HarborRequest represents the request message sent from Proxy to Dispatcher.
// It contains all the necessary information about the incoming HTTP request
// and a channel to send back the response.
type HarborRequest struct {

	// Context is the context for the request, used for cancellation and timeouts.
	Context context.Context

	// Url is the full URL of the incoming request, in string form.
	Url string

	// Method is the HTTP method of the incoming request.
	Method HTTPMethod

	// ContentType is the content type of the incoming request.
	// This is separated from Headers for easier access.
	ContentType ContentType

	// ResponseChannel is the channel through which the response will be sent
	// back to Proxy, allowing asynchronous communication.
	ResponseChannel chan HarborResponse

	// Headers contains the HTTP headers of the incoming request.
	// The headers are stored in lowercase to ensure case-insensitive access.
	Headers map[string]any

	// QueryParams contains the query parameters of the incoming request.
	QueryParams map[string]any

	// Body contains the body of the incoming request. It is stored as a map
	// and is generated directly from JSON,XML or plain structures.
	Body map[string]any
}
