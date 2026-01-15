// Package btp defines the Benthos–Tidal Protocol contract.
package btp

// ReferenceTransferData defines the envelope that groups all information returned
// to the client after a read operation for Minnow references.
type ReferenceTransferData struct {

	// References defines the set of Minnow references taken from read operation
	References []MinnowReference

	// Error defines the code related to an error occurred during read operation
	Error string

	// Size defines the number of elements included in [References] field
	Size int
}

// ProfileTransferData defines the envelope that groups all information returned
// to the client after a read operation for Minnow profiles.
type ProfileTransferData struct {

	// Profiles defines the set of Minnow references taken from read operation
	Profiles []MinnowProfile

	// Error defines the code related to an error occurred during read operation
	Error string

	// Size defines the number of elements included in [Profiles] field
	Size int
}

// EventTransferData defines the envelope that groups all information returned
// to the client after a read operation for Minnow events.
type EventTransferData struct {

	// Events defines the set of Minnow events taken from read operation
	Events []MinnowEvent

	// Error defines the code related to an error occurred during read operation
	Error string

	// Size defines the number of elements included in [Events] field
	Size int
}

// MinnowReference represents a reference to a Minnow profile stored in the search tree.
// It is used during the search and validation process to quickly identify the correct
// resource without including the full behavior definitions.
type MinnowReference struct {

	// Id defines the unique identifier for the resource.
	Id string

	// HttpMethod defines the HTTP method associated with the resource.
	Method string

	// Path defines the URL path for the resource.
	Path string

	// SpecialHeaders defines the special headers required for indexing the resource.
	SpecialHeaders map[string]string
}

// MinnowProfile represents a persisted MinnowProfile entity, corresponding to a saved reference
// in the registry. Each MinnowProfile defines an HTTP endpoint with its path, method,
// special headers, and associated behaviors. The struct includes metadata for
// tracking insertion and updates.
type MinnowProfile struct {

	// Behaviors defines the list of behaviors associated with this Minnow
	Behaviors []MinnowBehavior

	// Id defines the unique identifier of the Minnow, used for search indexing
	Id string

	// Method defines the HTTP method associated with the Minnow endpoint (GET, POST, etc.)
	Method string

	// URL path of the endpoint this Minnow represents
	Path string

	// SpecialHeaders defines the optional headers that uniquely distinguish
	// this reference or affect matching
	SpecialHeaders map[string]string
}

// MinnowBehavior represents a single behavior rule of a Minnow, defining a
// condition to match and the effect to execute if the condition is satisfied.
type MinnowBehavior struct {

	// Effect defines the effect to apply when the condition matches
	Effect MinnowBehaviorEffect

	// Condition defines the condition expression that triggers this behavior
	Condition string
}

// MinnowBehaviorEffect defines the action or response that occurs when a
// behavior is triggered. It can specify status code, headers, body content,
// content type, and optional execution logic.
type MinnowBehaviorEffect struct {

	// Type defines the type of effect to being applied
	Type string

	// Execute defines the script or command to execute when the behavior triggers
	Execute string

	// ContentType defines the MIME type of the response body
	ContentType string

	// RawBody defines the static raw response body as a string
	RawBody string

	// Headers defines the HTTP headers to include in the response
	Headers map[string]any

	// MappableBody defines the dynamic body map for templated responses
	MappableBody map[string]any

	// StatusCode defines the HTTP status code to return when the effect is executed
	StatusCode int
}

// MinnowEventType represents the type of edit made on Minnow profile.
// This is an alias for string to improve code readability and set predefined constants.
type MinnowEventType string

const (

	// MINNOW_EVENT_TYPE_ADD indicates that the edit operation refers to
	// adding a new Minnow profile.
	MINNOW_EVENT_TYPE_ADD MinnowEventType = "add"

	// MINNOW_EVENT_TYPE_EDIT_REFERENCE indicates that the edit operation refers
	// to adding an existing Minnow profile on index reference fields.
	MINNOW_EVENT_TYPE_ADD_REFERENCE MinnowEventType = "add_reference"

	// MINNOW_EVENT_TYPE_REMOVE indicates that the edit operation refers
	// to removing an existing Minnow profile.
	MINNOW_EVENT_TYPE_REMOVE MinnowEventType = "remove"

	// MINNOW_EVENT_TYPE_EDIT_REFERENCE indicates that the edit operation refers
	// to adding an existing Minnow profile on index reference fields.
	MINNOW_EVENT_TYPE_REMOVE_REFERENCE MinnowEventType = "remove_reference"

	// MINNOW_EVENT_TYPE_PROFILE_ACTIVATION indicates that the edit operation
	// refers to the activation status of a Minnow profile.
	MINNOW_EVENT_TYPE_PROFILE_ACTIVATION MinnowEventType = "profile_activation"

	// MINNOW_EVENT_TYPE_BEHAVIOR_ACTIVATION indicates that the edit operation
	// refers to the activation status of a Minnow behavior.
	MINNOW_EVENT_TYPE_BEHAVIOR_ACTIVATION MinnowEventType = "behavior_activation"

	// MINNOW_EVENT_TYPE_EDIT indicates that the edit operation refers to updating
	// an existing Minnow profile.
	MINNOW_EVENT_TYPE_EDIT MinnowEventType = "edit"
)

// MinnowEvent represents an edit operation performed on a Minnow entity.
// Each edit is stored as an event in a dedicated MongoDB collection, providing
// traceability and auditability of changes.
type MinnowEvent struct {

	// Operation contains metadata describing the type of edit and any previous reference.
	Operation MinnowEventOperation

	// MinnowId is the unique identifier of the Minnow that was edited.
	MinnowId string
}

// MinnowEventOperation defines the details of a specific edit operation
// applied to a Minnow, including its type and optional previous reference.
type MinnowEventOperation struct {

	// Reference holds the previous reference of the Minnow, if the operation
	// modifies its path, method, or headers. Nil if no previous reference exists.
	Reference MinnowEventReference

	// Type indicates the kind of edit performed (e.g., add, remove, update, add_reference).
	Type string
}

// MinnowEventReference represents the previous reference of a Minnow
// in case it was updated. It contains the HTTP method, path, and optional headers.
type MinnowEventReference struct {

	// Method is the HTTP method of the Minnow's previous reference.
	Method string

	// Path is the URL path of the Minnow's previous reference.
	Path string

	// SpecialHeaders defines any HTTP headers associated with the Minnow's previous reference.
	SpecialHeaders map[string]string
}
