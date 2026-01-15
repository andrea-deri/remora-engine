package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"remora/harbor/protocol/message"
	"remora/pkg/conversion"
	"remora/pkg/customerror"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// contextKey is a private alias type, used for defining key that can be
// included in the routine context
type contextKey string

const (

	// CONTEXT_KEY_REQUEST_ID is the context key used for unique request ID
	CONTEXT_KEY_REQUEST_ID contextKey = "requestID"
)

// RequestHandler is the HTTP request handler for the Harbor proxy. It implements the
// http.Handler interface by providing a ServeHTTP method to process incoming
// HTTP requests. The handler extracts relevant information from the request,
// such as headers, method, URL, query parameters, and body, and forwards
// this data to the dispatcher channel for further processing. It also handles
// response generation and timeout management.
type RequestHandler struct {

	// isReady defines the status flag about the engine readiness
	isReady atomic.Bool

	// BasePath defines the base path that will be excluded from URL extraction
	BasePath string

	// DispatcherChannel defines the channel on which the request for the second
	// level routing are sent in order to be handled
	DispatcherChannel chan message.HarborRequest

	// StatusChannel defines the dedicated channel where the status change messages
	// are sent in order to being processed by proxy handler
	StatusChannel chan string

	// RequestTimeout defines the time in seconds within which a request is considered
	// valid and after which it expires due to a timeout.
	RequestTimeout time.Duration
}

// Init initializes the RequestHandler context and starts the status change controller loop.
// This loop updates the readiness state in real-time, allowing readiness probes to reflect
// the current operational status of the handler.
func (handler *RequestHandler) Init() {

	handler.isReady.Store(false)
	log.Trace().Msgf("Initializing RequestHandler status with isReady: [%v]", handler.isReady.Load())

	// Launch a goroutine in order to handle status changes
	go handler.handleStatusChanges()
}

// handleStatusChanges listens for status commands from the StatusChannel and updates
// the RequestHandler's readiness state accordingly. Supported commands include "ready"
// (marks the handler as ready) and "unready" (marks the handler as not ready).
// This loop runs continuously until the StatusChannel is closed.
func (handler *RequestHandler) handleStatusChanges() {

	for statusCommand := range handler.StatusChannel {

		log.Trace().Msgf("A status command is arrived on RequestHandler: [%v]", statusCommand)

		switch statusCommand {

		// Handle status change when engine is ready to start processing
		case "ready":
			handler.isReady.Store(true)

		// Handle status change when engine is unready to process
		case "unready":
			handler.isReady.Store(false)
		}
	}
}

// ServeHTTP defines the procedure executed every time an HTTP request arrives on the listener.
// This implements the http.Handler interface and performs first-level routing: extracting
// information from the raw HTTP request and wrapping it into a HarborRequest for the
// second-level routing handled by the dispatcher.
func (handler *RequestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	// Block incoming request when engine is not fully ready
	if !handler.isReady.Load() {
		generateResponseForUnreadyEngine(writer)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), handler.RequestTimeout)
	ctx = context.WithValue(ctx, CONTEXT_KEY_REQUEST_ID, uuid.NewString())
	defer cancel()

	harborRequest := generateHarborRequest(ctx, request, handler)
	log.Debug().
		Str("RequestID", ctx.Value(CONTEXT_KEY_REQUEST_ID).(string)).
		Msgf("Received request [%s %s] with query parameters [%s], headers [%s], body [%s]", harborRequest.Method, harborRequest.Url, harborRequest.QueryParams, harborRequest.Headers, harborRequest.Body)

	handler.DispatcherChannel <- harborRequest

	select {

	// Channel has received computed response: send it in stream writer
	case response := <-harborRequest.ResponseChannel:
		writeValidResponse(ctx, writer, response, harborRequest)

	// Context has received a cancellation signal: send an error in stream writer
	case <-ctx.Done():
		writeTimeoutResponse(ctx, writer, harborRequest)
	}
}

// writeValidResponse writes a valid HTTP response produced by the dispatcher.
// The method sets standard headers, applies any custom response headers,
// writes the status code and encodes the body.
func writeValidResponse(ctx context.Context, writer http.ResponseWriter, response message.HarborResponse, harborRequest message.HarborRequest) {

	contentType := harborRequest.ContentType
	if contentType == "" {
		contentType = extractContentTypeFromString(response.ContentType)
	}

	responseHeaders := response.Headers
	setCustomStandardHeaders(writer, contentType)
	for key, value := range responseHeaders {
		if value != "" {
			writer.Header().Set(key, value)
		}
	}

	statusCode := response.StatusCode
	if statusCode == 0 {
		statusCode = 500
	}
	writer.WriteHeader(statusCode)

	log.Debug().
		Str("RequestID", ctx.Value(CONTEXT_KEY_REQUEST_ID).(string)).
		Msgf("Responding to request: [%s %s] with status code [%d] headers: [%s], body: [%s]", harborRequest.Method, harborRequest.Url, response.StatusCode, response.Headers, response.Body)

	encodeResponseBody(writer, response, contentType)
}

// writeTimeoutResponse writes a timeout response when the request context expires.
// The method builds a standard timeout HarborResponse, sets default headers,
// writes the status code and encodes the body.
func writeTimeoutResponse(ctx context.Context, writer http.ResponseWriter, harborRequest message.HarborRequest) {

	contentType := harborRequest.ContentType
	if contentType == "" {
		contentType = message.CONTENT_TYPE_JSON
	}

	// Building a response with custom error for timeout
	response := message.HarborResponse{
		Body:       message.BuildErrorResponse(customerror.CODE_HARBOR_REQUEST_TIMEOUT, "Request timeout"),
		StatusCode: http.StatusGatewayTimeout,
	}

	setCustomStandardHeaders(writer, contentType)
	writer.WriteHeader(response.StatusCode)

	log.Debug().
		Str("RequestID", ctx.Value(CONTEXT_KEY_REQUEST_ID).(string)).
		Msgf("Responding to request [%s %s] with status code [%d], headers [%s], body [%s]", harborRequest.Method, harborRequest.Url, response.StatusCode, response.Headers, response.Body)

	encodeResponseBody(writer, response, contentType)
}

// generateResponseForUnreadyEngine sends a standard 503 Service Unavailable response
// indicating that the REMORA Engine is not yet ready to accept requests.
// It sets standard headers, a Retry-After header and encodes the response body in JSON.
func generateResponseForUnreadyEngine(writer http.ResponseWriter) {

	log.Trace().Msgf("REMORA Engine is currently on readiness status [false]. Waiting for subprocesses to complete!")

	response := message.HarborResponse{
		Body:       message.BuildErrorResponse(customerror.CODE_ENGINE_UNREADY, "REMORA Engine is not yet ready to accept requests. Retry again after 10s!"),
		StatusCode: http.StatusServiceUnavailable,
		BodyType:   message.BODY_TYPE_MAP,
	}

	setCustomStandardHeaders(writer, message.CONTENT_TYPE_JSON)

	writer.Header().Set("retry-after", "10")
	writer.WriteHeader(response.StatusCode)

	encodeResponseBody(writer, response, message.CONTENT_TYPE_JSON)
}

// generateHarborRequest constructs a HarborRequest from an incoming HTTP request.
// It extracts headers, content type, HTTP method, URL (normalized to remove the handler's base path),
// query parameters, and body. A response channel is also created for sending back the processed response.
func generateHarborRequest(ctx context.Context, request *http.Request, handler *RequestHandler) message.HarborRequest {

	requestHeaders := extractHeaders(request)
	contentType := extractContentTypeFromString(requestHeaders["content-type"])
	httpMethod := extractMethod(request)
	url := extractURL(request, handler.BasePath)
	queryParams := extractQueryParameters(request)
	body := decodeRequestBody(request, contentType)

	return message.HarborRequest{
		Context:         ctx,
		Url:             url,
		Method:          httpMethod,
		Headers:         requestHeaders,
		ContentType:     contentType,
		QueryParams:     queryParams,
		Body:            body,
		ResponseChannel: make(chan message.HarborResponse),
	}
}

// decodeRequestBody extracts the request payload into a map for further processing.
// The decoding strategy depends on the Content-Type header: JSON and XML bodies are
// unmarshalled accordingly while unmapped content types are returned as a raw body map.
// Returns nil if the request has no body.
func decodeRequestBody(request *http.Request, contentType message.ContentType) map[string]any {

	body := request.Body
	if body == nil {
		// No body to decode: return instantly
		return nil
	}

	switch contentType {

	// Body in JSON format: extract it with JSON decoder
	case message.CONTENT_TYPE_JSON:
		var payload map[string]any
		error := json.NewDecoder(body).Decode(&payload)
		if !errors.Is(error, io.EOF) {
			// Ignoring EOF error: it means empty body
			log.Error().Str("Component", "Harbor").Msgf("Impossible to decode JSON body: %s", error)
		}
		return payload

	// Body in XML format: extract it with XML decoder
	case message.CONTENT_TYPE_XML:
		var payload map[string]any
		payload, error := conversion.DecodeFromXML(body)
		if error != nil && !errors.Is(error, io.EOF) {
			// Ignoring EOF error: it means empty body
			log.Error().Str("Component", "Harbor").Msgf("Impossible to decode XML body: %s", error)
		}
		return payload

	// Body in undefined format: extract it as plain content
	default:
		return map[string]any{"body": body}
	}
}

// encodeResponseBody writes the response payload to the HTTP response writer.
// The method checks the response's BodyType and the provided content type to determine
// the correct encoding strategy. Raw body content is written directly while map-based
// bodies are encoded according to the content type (JSON, XML or plain text).
// If the content type is unmapped, JSON encoding is used by default.
// No value is set if the body is nil.
func encodeResponseBody(writer http.ResponseWriter, response message.HarborResponse, contentType message.ContentType) {

	body := response.Body
	if body == nil {
		// No body to write: return instantly
		return
	}

	// Body as raw content: write it as a raw byte-array
	if response.BodyType == message.BODY_TYPE_RAW {
		_, err := writer.Write([]byte(response.Body.(string)))
		if err != nil {
			log.Error().
				Str("Component", "Harbor").
				Msgf("Impossible to write raw body on HTTP channel: %s", err)
		}
		return
	}

	switch contentType {

	// Body as mappable content, required as JSON: write it with JSON encoder
	case message.CONTENT_TYPE_JSON:
		err := json.NewEncoder(writer).Encode(response.Body)
		if err != nil {
			log.Error().
				Str("Component", "Harbor").
				Msgf("Impossible to write JSON body on HTTP channel: %s", err)
		}

	// Body as mappable content, required as plain text: write it as a raw byte-array
	case message.CONTENT_TYPE_STRING:
		_, err := writer.Write([]byte(response.Body.(string)))
		if err != nil {
			log.Error().
				Str("Component", "Harbor").
				Msgf("Impossible to write string body on HTTP channel: %s", err)
		}

	// Body as mappable content, required as XML: write it with XML encoder
	case message.CONTENT_TYPE_XML:
		xmlContent := conversion.EncodeAsXml(response.Body.(map[string]any))
		_, err := writer.Write(xmlContent)
		if err != nil {
			log.Error().
				Str("Component", "Harbor").
				Msgf("Impossible to write XML body on HTTP channel: %s", err)
		}

	// Body as mappable content, no format defined: write it with JSON encoder
	default:
		err := json.NewEncoder(writer).Encode(response.Body)
		if err != nil {
			log.Error().
				Str("Component", "Harbor").
				Msgf("Impossible to write generic body on HTTP channel: %s", err)
		}
	}

}

// setCustomStandardHeaders writes standard custom headers to the HTTP response.
// This includes a "x-powered-by" header for identification and a "Content-Type" header
// derived from the provided ContentType enum. Headers serve both as default values
// and as fallback if required headers are missing.
func setCustomStandardHeaders(writer http.ResponseWriter, contentType message.ContentType) {

	writer.Header().Set("x-powered-by", "REMORA Engine")
	writer.Header().Set("content-type", extractContentTypeFromEnum(contentType))
}
