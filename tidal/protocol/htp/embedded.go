// Package htp defines the Harbor–Tidal Protocol contract.
package htp

import (
	"remora/pkg/config"
	"remora/pkg/customerror"
	"remora/pkg/protocol/transfer/btp"
	"remora/pkg/protocol/transfer/htp"
	"remora/tidal/engine"

	"github.com/rs/zerolog/log"
)

// EmbeddedEngineRequestType represents the type of request required for Embedded Engine.
// This is an alias for string to improve code readability and set predefined constants.
// The accepted values are: and, or, clause.
type EmbeddedEngineRequestType int

const (

	// SEARCH_REQUEST defines the type of a required "search" request
	SEARCH_REQUEST EmbeddedEngineRequestType = 1

	// VALIDATION_REQUEST defines the type of a required "validation" request
	VALIDATION_REQUEST EmbeddedEngineRequestType = 2
)

// EmbeddedEngineRequest defines the structure of the message that is used to
// encapsulate requests to embedded implementation of Tidal Engine.
type EmbeddedEngineRequest struct {

	// Request represents the content that will be routed to required operation
	Request any

	// ReplyChannel represents the channel on which response content from Engine is sent
	ReplyChannel chan any

	// Type represents the enumerative value of the request type
	Type EmbeddedEngineRequestType
}

// EmbeddedEngine wraps a TidalEngine to provide an embedded execution environment.
// It exposes channels for sending operation requests and receiving status updates.
type EmbeddedEngine struct {

	// TidalEngine represents the concrete implementation of the Tidal Engine
	*engine.TidalEngine

	// RequestChannel represents the channel where the operation requests are sent
	RequestChannel chan EmbeddedEngineRequest

	// StatusChannel represents the channel on which the status changes are sent
	StatusChannel chan string
}

// NewEmbeddedEngine creates a new EmbeddedEngine instance wrapping the given TidalEngine.
// The returned engine is ready to be started, with channels initialized for requests and
// status notifications.
func NewEmbeddedEngine(engine *engine.TidalEngine) engine.Engine {

	return &EmbeddedEngine{
		TidalEngine:    engine,
		RequestChannel: make(chan EmbeddedEngineRequest),
	}
}

// Init initializes the EmbeddedEngine instance by initializing its internal TidalEngine
// and starting the operation controller loop. After initialization, the engine is ready
// to process requests in real time.
func (engine *EmbeddedEngine) Init(configMap config.ConfigMap, benthosAPI btp.BenthosAPI) {

	engine.TidalEngine.Init(configMap, benthosAPI)
	engine.StatusChannel <- "ready"
	go engine.StartListening()
}

// Search retrieves a Minnow reference from the search tree by using HTTP request information such
// as URL path, method and headers. This is the implementation exposed by the embedded Tidal Engine
// and provides asynchronous behavior for these operations.
func (engine *EmbeddedEngine) Search(request htp.SearchRequest) htp.SearchResult {

	replyChannel := make(chan any, 1)
	engine.RequestChannel <- EmbeddedEngineRequest{
		Request:      request,
		Type:         SEARCH_REQUEST,
		ReplyChannel: replyChannel,
	}
	rawResult := <-replyChannel
	return rawResult.(htp.SearchResult)
}

// Validate evaluates the input data against the behavior related to the Minnow referenced by
// the passed identifier and returns the resulting Effect result along with any error encountered.
// This is the implementation exposed by the embedded Tidal Engine and provides asynchronous
// behavior for these operations.
func (engine *EmbeddedEngine) Validate(request htp.ValidationRequest) htp.ValidationResult {

	replyChannel := make(chan any, 1)
	engine.RequestChannel <- EmbeddedEngineRequest{
		Request:      request,
		Type:         VALIDATION_REQUEST,
		ReplyChannel: replyChannel,
	}
	rawResult := <-replyChannel
	return rawResult.(htp.ValidationResult)
}

// StartListening begins processing requests sent to the EmbeddedEngine's request channel.
// Each request represents an operation invoked by a client (e.g., Harbor Dispatcher), such as
// Search or Validate. Requests are handled asynchronously, and the method continuously listens
// for incoming messages, dispatching them to the appropriate handler based on the request type.
func (engine *EmbeddedEngine) StartListening() {

	for request := range engine.RequestChannel {

		requestType := request.Type
		switch requestType {

		case SEARCH_REQUEST:
			engine.executeSearch(request)

		case VALIDATION_REQUEST:
			engine.executeValidate(request)

		default:
			panic("An error occurred while receiving request on embedded Tidal engine. Invalid request type.")
		}
	}
}

// executeSearch handles a SearchRequest received on the EmbeddedEngine's request channel.
// It attempts to cast the incoming request to a SearchRequest type. If successful, it delegates
// the search to the underlying TidalEngine. The result, or a fallback indicating failure if
// casting fails, is returned via the request's reply channel.
func (engine *EmbeddedEngine) executeSearch(request EmbeddedEngineRequest) {

	var searchResult htp.SearchResult
	searchRequest, isConversionOk := request.Request.(htp.SearchRequest)

	if !isConversionOk {
		log.Error().Msgf("An error occurred while receiving search request. Cast to SearchRequest type failed.")
		searchResult = htp.SearchResult{IsFound: false}
	} else {
		searchResult = engine.TidalEngine.Search(searchRequest)
	}

	request.ReplyChannel <- searchResult
}

// executeValidate handles a ValidationRequest received on the EmbeddedEngine's request channel.
// It attempts to cast the incoming request to a ValidationRequest type. If successful, it
// delegates the validation to the underlying TidalEngine. The result, or an error if casting
// fails, is sent back via the request's reply channel.
func (engine *EmbeddedEngine) executeValidate(request EmbeddedEngineRequest) {

	var validationResult htp.ValidationResult
	validationRequest, isConversionOk := request.Request.(htp.ValidationRequest)

	if !isConversionOk {
		log.Error().Msgf("An error occurred while receiving validation request. Cast to ValidationRequest type failed.")
		validationResult = htp.ValidationResult{Error: string(customerror.CODE_TIDAL_VALIDATION_INVALID_REQUEST)}
	} else {
		validationResult = engine.TidalEngine.Validate(validationRequest)
	}

	request.ReplyChannel <- validationResult
}
