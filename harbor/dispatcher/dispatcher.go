package dispatcher

import (
	"fmt"
	"os"
	"strings"

	"remora/harbor/protocol/message"
	"remora/pkg/config"
	"remora/pkg/conversion/format"
	"remora/pkg/conversion/maps"
	"remora/pkg/customerror"
	"remora/pkg/protocol/transfer/htp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// HarborDispatcher manages incoming requests, coordinating with the Tidal Engine.
// It handles request queuing, concurrency control, and response generation.
type HarborDispatcher struct {

	// logger is the zerolog logger instance for logging within the dispatcher.
	logger zerolog.Logger

	// specialEndpointPrefix defines the prefix path on which the "special operation"
	// endpoints are included. If not set, it is defined with "/r" path.
	specialEndpointPrefix string

	// LogLevel defines the logging level for the dispatcher component.
	LogLevel string

	// tidalEngine define the interface for the communication with Tidal Engine.
	// It is used in order to search and validate input data against Minnow Profiles.
	tidalEngine htp.TidalAPI

	// InputChannel is the channel through which incoming requests are received.
	// Requests are sent to this channel by the HarborProxy component.
	//
	// The channel has a buffer size defined by MaxAcceptedRequests. If the channel
	// is full, new incoming requests will be rejected by Dispatcher. Requests
	// received on this channel are processed concurrently, up to the limit
	// defined by MaxConcurrentRequests.
	InputChannel chan message.HarborRequest

	// requestsSemaphore is a semaphore channel to limit concurrent request processing.
	// It has a buffer size defined by MaxConcurrentRequests and each request processing
	// goroutine must acquire a slot in the semaphore before starting processing, and
	// release it when done.
	requestsSemaphore chan struct{}

	// MaxConcurrentRequests defines the maximum number of requests
	// that can be processed concurrently by the engine. If the limit
	// is reached, new incoming requests will be queued until a processing
	// slot is available.
	//
	// This is a hard limit: requests are executed only when a processing
	// slot is available in the semaphore channel.
	MaxConcurrentRequests int

	// MaxAcceptedRequests defines the maximum number of accepted requests
	// that can be queued in the input channel. If the channel is full,
	// new incoming requests will be rejected by Dispatcher.
	//
	// This is a soft limit: requests are accepted until the channel
	// is full. After that, new requests will be blocked by Dispatcher
	// returning a 429 HTTP code to the client.
	MaxAcceptedRequests int

	// CanCacheContent defines if the dispatcher can cache content internally
	// to improve performance on repeated requests.
	//
	// Note: This feature is not implemented yet.
	CanCacheContent bool
}

// NewDispatcher constructs and returns a fully configured HarborDispatcher instance.
// Configuration values are read from the provided ConfigMap, applying defaults when keys are missing.
// The returned dispatcher is ready to be started by the caller.
func NewDispatcher(configMap config.ConfigMap, tidalAPIs htp.TidalAPI) HarborDispatcher {

	specialEndpointPrefix := configMap.ReadString("DISPATCHER_SPECIAL_ENDPOINT_PREFIX", "r")
	specialEndpointPrefix = fmt.Sprintf("%s/", strings.Trim(specialEndpointPrefix, "/"))

	return HarborDispatcher{
		MaxAcceptedRequests:   configMap.ReadInt("DISPATCHER_MAX_ACCEPTED_REQUESTS", 64),
		MaxConcurrentRequests: configMap.ReadInt("MAX_CONCURRENT_REQUESTS", 32),
		CanCacheContent:       configMap.ReadBoolean("CAN_CACHE_CONTENT", false),
		tidalEngine:           tidalAPIs,
		specialEndpointPrefix: specialEndpointPrefix,
		LogLevel:              configMap.ReadString("HARBOR_DISPATCHER_LOGGING_LEVEL", "info"),
	}
}

// Start prepares the dispatcher for handling requests by initializing channels
// and concurrency control. After setup, the processing loop is launched in a separate
// goroutine and begins consuming incoming requests.
func (dispatcher *HarborDispatcher) Start() {

	// Initializing the logger instance with defined level
	level, err := zerolog.ParseLevel(dispatcher.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	dispatcher.logger = zerolog.New(os.Stdout).
		With().
		Str("Component", "Harbor").
		Str("Module", "Dispatcher").
		Logger().
		Level(level)

	dispatcher.InputChannel = make(chan message.HarborRequest, dispatcher.MaxAcceptedRequests)
	dispatcher.requestsSemaphore = make(chan struct{}, dispatcher.MaxConcurrentRequests)

	log.Info().Msgf("Starting dispatcher...")
	go dispatcher.run()
}

// run is the Dispatcher’s main processing loop. It consumes incoming requests
// from InputChannel and delegates each to a dedicated goroutine. Concurrency is
// bounded by requestsSemaphore, which ensures that active request handlers never
// exceed MaxConcurrentRequests. The method processes requests in arrival order
// and continues running until InputChannel is closed.
func (dispatcher *HarborDispatcher) run() {

	log.Info().Msgf("Dispatcher module is running! Max accepted requests: [%d], Max concurrent requests: [%d]", dispatcher.MaxAcceptedRequests, dispatcher.MaxConcurrentRequests)

	// Each request is processed as the channel delivers it. Backpressure occurs
	// naturally if InputChannel reaches capacity (handled upstream by HarborProxy),
	// while concurrency is constrained by acquiring a slot from requestsSemaphore.
	for request := range dispatcher.InputChannel {

		// Acquire a concurrency slot before spawning the processing goroutine.
		select {

		case dispatcher.requestsSemaphore <- struct{}{}:
			log.Trace().Msg("Semaphore slot acquired correctly")

		case <-request.Context.Done():
			log.Debug().Msg("Request expired before acquiring semaphore")
			continue
		}

		// Spawn a single goroutine that handle one of the incoming request.
		// Each goroutine will release its concurrency slot when done.
		log.Trace().Msgf("Dispatched request: %s", request.Body)
		go dispatcher.handleRequest(request)
	}
}

// handleRequest processes a single HarborRequest. It determines whether the request
// targets a special or standard endpoint, delegates handling accordingly and then
// attempts to return a HarborResponse through the request’s ResponseChannel.
// The send operation is guarded in order to avoid blocking if the request context
// has expired. The associated semaphore slot is released by the caller once execution
// completes.
func (dispatcher *HarborDispatcher) handleRequest(request message.HarborRequest) {

	// Ensure concurrency slot is always released.
	defer func() {
		recoveredState := recover()
		if recoveredState != nil {
			log.Error().Msgf("Recovered panic for handleRequest in Dispatcher. Releasing request semaphore...")
		}
		<-dispatcher.requestsSemaphore
	}()

	var response message.HarborResponse
	if isSpecialEndpoint(dispatcher.specialEndpointPrefix, request) {

		response = dispatcher.handleSpecialRequest(request)

	} else {

		response = dispatcher.handleStandardRequest(request)
	}

	select {

	// Send succeeds only if the receiving side is still active. This operation is non-blocking
	// in order to prevent hanging if the Proxy is no longer listening (e.g. due to client
	// disconnection or socket timeout).
	case request.ResponseChannel <- response:
		log.Debug().Msgf("Sent response from Dispatcher: %v", response)

	// If the request context is canceled, skip the send to avoid stale or blocked writes.
	case <-request.Context.Done():
		log.Warn().Msgf("Skipped sending response for canceled request: %s", request.Body)
	}
}

// handleSpecialRequest processes internal REMORA Engine operations exposed through
// special endpoints. The endpoint prefix is removed before routing so only the
// core operation name is evaluated.
//
// Unsupported operations return a structured error response with status 400.
func (dispatcher *HarborDispatcher) handleSpecialRequest(request message.HarborRequest) message.HarborResponse {

	// Sanitize endpoint path in order to exclude prefix, not used in 2nd-level routing
	endpointPath := strings.Replace(request.Url, dispatcher.specialEndpointPrefix, "", 1)

	response := message.HarborResponse{}
	switch endpointPath {

	// Health check endpoint
	case "ping":
		response.StatusCode = 200
		response.BodyType = message.BODY_TYPE_MAP

	// Force refresh endpoint
	case "refresh":
		response.StatusCode = 400
		response.BodyType = message.BODY_TYPE_MAP
		response.Body = message.BuildErrorResponse(customerror.CODE_ENGINE_UNIMPLEMENTED_FEATURE, "Internal API [refresh] not implemented yet!")

	// Unmapped API: return a bad request
	default:
		response.StatusCode = 400
		response.BodyType = message.BODY_TYPE_MAP
		response.Body = message.BuildErrorResponse(customerror.CODE_ENGINE_UNIMPLEMENTED_FEATURE, fmt.Sprintf("Internal API [%v] not implemented yet!", endpointPath))
	}

	return response
}

// handleStandardRequest processes standard REMORA engine requests, i.e. requests
// that map to simulated Minnow behavior. It performs resource lookup through and,
// if found, execute validation and response generation.
// The concurrency slot acquired by the dispatcher loop is released upon completion.
func (dispatcher *HarborDispatcher) handleStandardRequest(request message.HarborRequest) message.HarborResponse {

	searchRequest := htp.SearchRequest{
		Ctx:     request.Context,
		Path:    request.Url,
		Method:  string(request.Method),
		Headers: request.Headers,
		Input:   make(map[string]any),
	}

	// Perform resource lookup via the Tidal Engine.
	searchResult := dispatcher.tidalEngine.Search(searchRequest)

	// If no resource matches the incoming request, return a custom error with HTTP code 404
	if !searchResult.IsFound {
		responseMessage := fmt.Sprintf("No valid resource is found at path [%s %s]", request.Method, request.Url)
		return generateErrorResponse(404, responseMessage, customerror.CODE_HARBOR_RESOURCE_NOT_FOUND)
	}

	searchRequest.Input["pparam"] = searchResult.PathParams
	searchRequest.Input["header"] = request.Headers
	searchRequest.Input["qparam"] = request.QueryParams
	searchRequest.Input["body"] = request.Body

	// Convert searchRequest in lowercase in order to avoid mismatches during
	// validation caused by text with different case
	searchRequest.Input = maps.ToLowerKeys(searchRequest.Input)

	validationRequest := htp.ValidationRequest{
		MinnowId: searchResult.ResourceId,
		Input:    searchRequest.Input,
	}
	return dispatcher.validateRequest(validationRequest, request)
}

// validateRequest performs behavior validation for a matched resource using the Tidal Engine.
// It evaluates whether the resource defines a valid behavior, reports any validation errors and
// constructs the appropriate HarborResponse based on the resulting content definition.
// Validation failure yields standardized error responses.
func (dispatcher *HarborDispatcher) validateRequest(validationRequest htp.ValidationRequest, request message.HarborRequest) message.HarborResponse {

	validationResult := dispatcher.tidalEngine.Validate(validationRequest)
	effectResult := validationResult.Result

	var response message.HarborResponse

	// No valid behavior associated with the resource.
	if !effectResult.IsValid {
		responseMessage := fmt.Sprintf("An error occurred during behavior validation. No valid behavior is found at path [%s %s]", request.Method, request.Url)
		return generateErrorResponse(404, responseMessage, customerror.CODE_HARBOR_RESOURCE_WITHOUT_BEHAVIOR)
	}

	// Validation engine reported an internal error while evaluating behavior rules.
	if validationResult.Error != "" {
		responseMessage := fmt.Sprintf("An error occurred during behavior validation. [%s]", validationResult.Error)
		return generateErrorResponse(500, responseMessage, customerror.CODE_HARBOR_RESOURCE_WITH_INVALID_BEHAVIOR)
	}

	if effectResult.RawBody != "" {
		response.Body = effectResult.RawBody
		response.BodyType = message.BODY_TYPE_RAW
	} else {
		response.Body = effectResult.MappableBody
		response.BodyType = message.BODY_TYPE_MAP
	}

	response.StatusCode = effectResult.StatusCode
	response.ContentType = effectResult.ContentType
	response.Headers = format.FromAnyToStringMap(effectResult.Headers)

	return response
}

// generateErrorResponse constructs a standardized HarborResponse representing an error.
// It logs the message, sets the HTTP status code and returns a body containing the
// structured error payload built with the given error code.
func generateErrorResponse(statusCode int, msg string, errorCode customerror.RemoraErrorCode) message.HarborResponse {

	log.Error().Msgf("%s", msg)
	return message.HarborResponse{
		StatusCode: statusCode,
		BodyType:   message.BODY_TYPE_MAP,
		Body:       message.BuildErrorResponse(errorCode, msg),
	}
}

// isSpecialEndpoint determines whether the request targets a special internal endpoint.
// Special endpoints use a distinct URL prefix reserved for internal engine operations.
func isSpecialEndpoint(specialEndpointPrefix string, request message.HarborRequest) bool {

	return strings.HasPrefix(request.Url, specialEndpointPrefix)
}
