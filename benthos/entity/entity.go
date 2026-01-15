package entity

// MinnowProfile represents a persisted Minnow Profile entity saved in the registry.
// Each MinnowProfile defines an HTTP endpoint with its path, method, special headers,
// and associated behaviors. The struct includes metadata for tracking insertion and updates.
type MinnowProfile struct {

	// Behaviors defines the list of behaviors associated with this Minnow
	Behaviors []MinnowBehavior `bson:"behaviors"`

	// Id defines the unique identifier of the Minnow, used for search indexing
	Id string `bson:"id"`

	// Method defines the HTTP method associated with the Minnow endpoint (GET, POST, etc.)
	Method string `bson:"method"`

	// URL path of the endpoint this Minnow represents
	Path string `bson:"path"`

	// CorrelationId defines the correlation ID used for grouping Minnows
	CorrelationId string `bson:"correlation_id"`

	// SpecialHeaders defines the optional headers that uniquely distinguish
	// this reference or affect matching
	SpecialHeaders map[string]string `bson:"special_headers"`

	// InsertedAt defines the timestamp when the Minnow was persisted
	InsertedAt int64 `bson:"inserted_at"`

	// UpdatedAt defines the timestamp of the last update to the Minnow
	UpdatedAt int64 `bson:"updated_at"`
}

// MinnowBehavior represents a single behavior strategy on a Minnow Profile, defining a
// condition to match and the effect to apply if the condition is satisfied.
type MinnowBehavior struct {

	// Effect defines the effect to apply when the condition matches
	Effect MinnowBehaviorEffect `bson:"effect"`

	// Id defines the unique identifier of the behavior
	Id string `bson:"id"`

	// Condition defines the condition expression that triggers this behavior
	Condition string `bson:"condition"`

	// IsActive defines whether the behavior is active and should be evaluated
	IsActive bool `bson:"is_active"`
}

// MinnowBehaviorEffect defines the action or response that occurs when a
// behavior is triggered. It can specify status code, headers, body content,
// content type and optional execution logic.
type MinnowBehaviorEffect struct {

	// Type defines the type of effect to being applied
	Type string `bson:"type"`

	// ContentType defines the MIME type of the response body
	ContentType string `bson:"content_type"`

	// RawBody defines the static raw response body as a string
	RawBody string `bson:"raw_body"`

	// Execute defines the script or command to execute when the behavior triggers
	Execute string `bson:"execute"`

	// Headers defines the HTTP headers to include in the response
	Headers map[string]any `bson:"headers"`

	// MappableBody defines the dynamic body map for templated responses
	MappableBody map[string]any `bson:"mappable_body"`

	// StatusCode defines the HTTP status code to return when the effect is executed
	StatusCode int `bson:"status_code"`
}

// MinnowEvent represents an operation performed on a Minnow profile entity.
// Each operation is stored as an atomic event, providing traceability and auditability
// of changes.
type MinnowEvent struct {

	// Operation contains metadata describing the type of edit and any previous reference.
	Operation MinnowEventOperation `bson:"operation"`

	// Context holds information about the origin and timing of the event.
	Context MinnowEventContext `bson:"context"`

	// EventId is the unique identifier of the event in the database.
	EventId string `bson:"_id,omitempty"`

	// MinnowId is the unique identifier of the Minnow that was edited.
	MinnowId string `bson:"minnow_id"`
}

// MinnowEventOperation defines the details of a specific operation applied to a Minnow
// Profile, including its type and optional previous reference.
type MinnowEventOperation struct {

	// Reference holds the previous reference of the Minnow, if the operation
	// modifies its path, method, or headers.
	// It is defined as nil if no previous reference exists.
	Reference MinnowEventReference `bson:"reference,omitempty"`

	// Type indicates the kind of edit performed (e.g., add, remove, update, add_reference).
	Type string `bson:"type"`
}

// MinnowEventReference represents the previous reference of a Minnow Profile
// in case it was updated. It contains the HTTP method, path, and optional headers.
type MinnowEventReference struct {

	// Method is the HTTP method of the Minnow's previous reference.
	Method string `bson:"method"`

	// Path is the URL path of the Minnow's previous reference.
	Path string `bson:"path"`

	// SpecialHeaders defines any HTTP headers associated with the Minnow's previous reference.
	SpecialHeaders map[string]string `bson:"special_headers,omitempty"`
}

// MinnowEventContext stores metadata about the source and timing of a Minnow edit event.
type MinnowEventContext struct {

	// Source identifies the user, service, or system that triggered the edit.
	Source string `bson:"source"`

	// Timestamp records the exact time when the edit operation was executed.
	Timestamp int64 `bson:"timestamp"`
}
