package proxy

import (
	"net/http"
	"remora/harbor/protocol/message"
	"strings"
)

// extractQueryParameters normalizes query parameters into a lower-cased string map.
// Each key is converted to lower case. Multiple values for the same parameter are
// joined with commas. The returned map contains raw, unmodified values.
func extractQueryParameters(request *http.Request) map[string]any {

	queryParameters := make(map[string]any)

	rawQueryParameters := request.URL.Query()
	for key, value := range rawQueryParameters {
		lowercaseKey := strings.ToLower(key)
		queryParameters[lowercaseKey] = strings.Join(value, ",")
	}

	return queryParameters
}

// extractHeaders normalizes HTTP headers into a lower-cased string map.
// Each key is converted to lower case. Multiple values for the same header
// are joined with commas. The returned map contains raw, unmodified values.
func extractHeaders(request *http.Request) map[string]any {

	headers := make(map[string]any)

	rawHeaders := request.Header
	for key, value := range rawHeaders {
		lowercaseKey := strings.ToLower(key)
		headers[lowercaseKey] = strings.Join(value, ",")
	}

	return headers
}

// extractContentTypeFromString maps a raw Content-Type header value to an internal
// ContentType enum. Matching is case-insensitive and performed via substring search.
// If no known type is detected, an empty ContentType is returned.
func extractContentTypeFromString(headerRawValue any) message.ContentType {

	if headerRawValue == nil {
		return ""
	}

	// The header value is converted in lower-case, avoiding possible errors in comparation
	headerValue := headerRawValue.(string)
	headerValue = strings.ToLower(headerValue)

	switch {

	case strings.Contains(headerValue, string(message.CONTENT_TYPE_JSON)):
		return message.CONTENT_TYPE_JSON

	case strings.Contains(headerValue, string(message.CONTENT_TYPE_XML)):
		return message.CONTENT_TYPE_XML
	}

	return ""
}

// extractContentTypeFromEnum converts an internal ContentType enum value into a
// concrete MIME type string. If the enum is unknown, "text/plain" is returned.
func extractContentTypeFromEnum(headerValue message.ContentType) string {

	switch headerValue {

	case message.CONTENT_TYPE_JSON:
		return "application/json"

	case message.CONTENT_TYPE_XML:
		return "text/xml"
	}

	return "text/plain"
}

// extractMethod converts the HTTP method from a request into an enumerated HTTPMethod value.
// Unrecognized methods default to HTTP_METHOD_GET.
func extractMethod(request *http.Request) message.HTTPMethod {

	method := strings.ToLower(request.Method)
	switch method {

	case "get":
		return message.HTTP_METHOD_GET
	case "post":
		return message.HTTP_METHOD_POST
	case "put":
		return message.HTTP_METHOD_PUT
	case "delete":
		return message.HTTP_METHOD_DELETE
	case "head":
		return message.HTTP_METHOD_HEAD
	case "options":
		return message.HTTP_METHOD_OPTIONS
	case "patch":
		return message.HTTP_METHOD_PATCH
	case "trace":
		return message.HTTP_METHOD_TRACE
	case "connect":
		return message.HTTP_METHOD_CONNECT
	}

	return message.HTTP_METHOD_GET
}

// extractURL returns the URL path from the HTTP request, optionally removing a specified base path.
func extractURL(request *http.Request, basePath string) string {

	urlPath := request.URL.Path
	if basePath != "" {
		urlPath = strings.TrimPrefix(urlPath, basePath)
	}
	return urlPath
}
