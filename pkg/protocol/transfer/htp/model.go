// Package htp defines the Harbor–Tidal Protocol contract.
package htp

import (
	"context"
	"remora/pkg/coral/profile"
)

// SearchRequest represents the request of a search operation.
type SearchRequest struct {

	// Ctx is the context for the search operation, used for cancellation and deadlines.
	Ctx context.Context

	// Path is the URL path to search for.
	Path string

	// Method is the HTTP method to use in the search.
	Method string

	// Headers contains the HTTP headers for the search request.
	Headers map[string]any

	// Input is the dynamic input map for search parameters.
	Input map[string]any
}

// SearchResult represents the result of a search operation made on the search tree.
type SearchResult struct {

	// ResourceId is the unique identifier of the found resource.
	ResourceId string

	// PathParams contains the path parameters extracted from wildcard tokens.
	PathParams map[string]any

	// IsFound indicates if the search found a valid result.
	IsFound bool
}

// ValidationRequest represents the request of a validation operation.
type ValidationRequest struct {

	// MinnowID is the identifier of the Minnow to be used during validation.
	MinnowId string

	// Input is the dynamic input map used as parameters for validation.
	Input map[string]any
}

// ValidationResult represents the result of a validation operation made on the
// retrieved Minnow profile.
type ValidationResult struct {

	// Result defines the resulting effect of the validation process, generated
	// from the execution of the Minnow behavior.
	Result profile.EffectResult

	// Error defines the code related to an error occurred during validation
	Error string
}
