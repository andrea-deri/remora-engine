// Package htp defines the Harbor–Tidal Protocol contract.
package htp

// TidalAPI defines the guidelines on which all the Tidal implementations
// must follows in order to exposes the required behaviors on rule engine.
type TidalAPI interface {

	// Search retrieves a Minnow reference from the search tree by
	// using HTTP request information such as URL path, method, and headers.
	Search(request SearchRequest) SearchResult

	// Validate evaluates the input data against the behavior related to the
	// Minnow referenced by the passed identifier and returns the resulting
	// Effect result along with any error encountered.
	Validate(request ValidationRequest) ValidationResult
}
