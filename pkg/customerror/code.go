package customerror

// RemoraErrorCode represents the type of the custom codes generated in case of error.
// This is an alias for string to improve code readability and set predefined constants.
type RemoraErrorCode string

const (

	// CODE_ENGINE_UNREADY is used when REMORA engine is not ready to accept requests yet
	CODE_ENGINE_UNREADY RemoraErrorCode = "00100"

	// CODE_ENGINE_UNIMPLEMENTED_FEATURE is used when REMORA engine has execute some unimplemented feature.
	// Use this if you want to return an error and not panic the application.
	CODE_ENGINE_UNIMPLEMENTED_FEATURE RemoraErrorCode = "00199"

	// ------------------

	// CODE_HARBOR_RESOURCE_NOT_FOUND is used when no valid resource is found at request path
	CODE_HARBOR_RESOURCE_NOT_FOUND RemoraErrorCode = "01000"

	// CODE_HARBOR_REQUEST_TIMEOUT is used when Harbor context sends termination signal due to timeout
	CODE_HARBOR_REQUEST_TIMEOUT RemoraErrorCode = "01001"

	// CODE_HARBOR_RESOURCE_WITHOUT_BEHAVIOR is used when a resource is found but hasn't a behavior
	CODE_HARBOR_RESOURCE_WITHOUT_BEHAVIOR RemoraErrorCode = "01002"

	// CODE_HARBOR_RESOURCE_WITH_INVALID_BEHAVIOR is used when a resource is found but hasn't a valid behavior
	CODE_HARBOR_RESOURCE_WITH_INVALID_BEHAVIOR RemoraErrorCode = "01003"

	// CODE_HARBOR_GZIP_NOT_DECODABLE is used when Harbor context sends termination signal due to GZIP decode error
	CODE_HARBOR_GZIP_NOT_DECODABLE RemoraErrorCode = "01004"

	// ------------------

	// CODE_CORAL_SYNTAXTREE_UNMAPPED_CLAUSE is used when a clause is unmapped in syntax tree
	CODE_CORAL_SYNTAXTREE_UNMAPPED_CLAUSE RemoraErrorCode = "02000"

	// CODE_CORAL_SYNTAXTREE_OPERATOR_NOT_IMPLEMENTED is used when an operator is not implemented in syntax tree
	CODE_CORAL_SYNTAXTREE_OPERATOR_NOT_IMPLEMENTED RemoraErrorCode = "02001"

	// CODE_CORAL_SYNTAXTREE_FUNCTIONALITY_NOT_IMPLEMENTED is used when a functionality is not implemented in syntax tree
	CODE_CORAL_SYNTAXTREE_FUNCTIONALITY_NOT_IMPLEMENTED RemoraErrorCode = "02002"

	// CODE_CORAL_SYNTAXTREE_REGEX_NOT_APPLICABLE is used when regex is not applicable in syntax tree
	CODE_CORAL_SYNTAXTREE_REGEX_NOT_APPLICABLE RemoraErrorCode = "02003"

	// CODE_CORAL_SYNTAXTREE_MALFORMED_LEXEME is used when a lexeme is malformed in syntax tree
	CODE_CORAL_SYNTAXTREE_MALFORMED_LEXEME RemoraErrorCode = "02004"

	// CODE_CORAL_SYNTAXTREE_INVALID_STARTING_TOKEN_LEXEME is used when a lexeme starts with an invalid token in syntax tree
	CODE_CORAL_SYNTAXTREE_INVALID_STARTING_TOKEN_LEXEME RemoraErrorCode = "02005"

	// CODE_CORAL_SYNTAXTREE_INVALID_ENDING_TOKEN_LEXEME is used when a lexeme ends with an invalid token in syntax tree
	CODE_CORAL_SYNTAXTREE_INVALID_ENDING_TOKEN_LEXEME RemoraErrorCode = "02006"

	// CODE_CORAL_SYNTAXTREE_MISSING_RIGHT_PARENTHESYS is used when a right parenthesis is missing in syntax tree
	CODE_CORAL_SYNTAXTREE_MISSING_RIGHT_PARENTHESYS RemoraErrorCode = "02007"

	// CODE_CORAL_SYNTAXTREE_MISSING_LITERAL_TOKEN is used when a literal token is missing in syntax tree
	CODE_CORAL_SYNTAXTREE_MISSING_LITERAL_TOKEN RemoraErrorCode = "02008"

	// CODE_CORAL_SYNTAXTREE_INVALID_OPERATOR_TOKEN is used when an operator token is invalid in syntax tree
	CODE_CORAL_SYNTAXTREE_INVALID_OPERATOR_TOKEN RemoraErrorCode = "02009"

	// CODE_CORAL_SYNTAXTREE_INVALID_RIGHT_OPERATOR_TOKEN is used when a right operator token is invalid in syntax tree
	CODE_CORAL_SYNTAXTREE_INVALID_RIGHT_OPERATOR_TOKEN RemoraErrorCode = "02010"

	// CODE_CORAL_SYNTAXTREE_INVALID_REGEX_OPERATOR_TOKEN is used when a regex operator token is invalid in syntax tree
	CODE_CORAL_SYNTAXTREE_INVALID_REGEX_OPERATOR_TOKEN RemoraErrorCode = "02011"

	// CODE_CORAL_SYNTAXTREE_MISSING_THEN_TOKEN is used when a THEN token is missing in syntax tree
	CODE_CORAL_SYNTAXTREE_MISSING_THEN_TOKEN RemoraErrorCode = "02012"

	// CODE_CORAL_SYNTAXTREE_MISSING_RETURN_TOKEN is used when a RETURN token is missing in syntax tree
	CODE_CORAL_SYNTAXTREE_MISSING_RETURN_TOKEN RemoraErrorCode = "02013"

	// CODE_CORAL_SYNTAXTREE_MISSING_EXEC_LITERAL_TOKEN is used when an EXEC literal token is missing in syntax tree
	CODE_CORAL_SYNTAXTREE_MISSING_EXEC_LITERAL_TOKEN RemoraErrorCode = "02014"

	// CODE_CORAL_SYNTAXTREE_MISSING_RETURN_LITERAL_TOKEN is used when a RETURN literal token is missing in syntax tree
	CODE_CORAL_SYNTAXTREE_MISSING_RETURN_LITERAL_TOKEN RemoraErrorCode = "02015"

	// CODE_CORAL_OPERATOR_INVALID_BOOLEAN_CONVERSION is used when boolean conversion fails in operator evaluation
	CODE_CORAL_OPERATOR_INVALID_BOOLEAN_CONVERSION RemoraErrorCode = "02016"

	// CODE_CORAL_OPERATOR_INVALID_STRING_CONVERSION is used when string conversion fails in operator evaluation
	CODE_CORAL_OPERATOR_INVALID_STRING_CONVERSION RemoraErrorCode = "02017"

	// CODE_CORAL_OPERATOR_INVALID_NUMERIC_CONVERSION is used when numeric conversion fails in operator evaluation
	CODE_CORAL_OPERATOR_INVALID_NUMERIC_CONVERSION RemoraErrorCode = "02018"

	// ------------------

	// CODE_CORAL_RUNNER_UNMARSHAL_FAILED is used when unmarshalling content in Msgpack format fails
	CODE_CORAL_RUNNER_UNMARSHAL_FAILED RemoraErrorCode = "02100"

	// CODE_CORAL_RUNNER_RAWBODY_RESOLUTION_FAILED is used when value resolution on raw body fails
	CODE_CORAL_RUNNER_RAWBODY_RESOLUTION_FAILED RemoraErrorCode = "02101"

	// CODE_CORAL_RUNNER_MAPPABLEBODY_RESOLUTION_FAILED is used when value resolution on mappable body fails
	CODE_CORAL_RUNNER_MAPPABLEBODY_RESOLUTION_FAILED RemoraErrorCode = "02102"

	// CODE_CORAL_RUNNER_MAPPABLEBODY_RESOLUTION_FAILED is used when value resolution on headers fails
	CODE_CORAL_RUNNER_HEADERS_RESOLUTION_FAILED RemoraErrorCode = "02103"

	// ------------------

	// CODE_TIDAL_SEARCH_PROFILE_NOT_FOUND is used when Tidal search cannot find a valid profile
	// in profile map with passed data.
	CODE_TIDAL_SEARCH_PROFILE_NOT_FOUND RemoraErrorCode = "03000"

	// ------------------

	// CODE_TIDAL_VALIDATION_INVALID_REQUEST is used when Tidal validation cannot starts on request
	// because the internal wrapped message contains an invalid request message.
	CODE_TIDAL_VALIDATION_INVALID_REQUEST RemoraErrorCode = "03100"

	// ------------------

	// CODE_BENTHOS_CONNECTION_FAILED is used when Benthos cannot instantiate a connection with
	// persistence provider.
	CODE_BENTHOS_CONNECTION_FAILED RemoraErrorCode = "04000"

	// CODE_BENTHOS_ENTITY_NOT_FOUND is used when Benthos cannot finalize a "read"
	// operation on Profiles for MongoDB provider because no valid entity is found.
	CODE_BENTHOS_ENTITY_NOT_FOUND RemoraErrorCode = "04001"

	// CODE_BENTHOS_MONGODB_PROFILES_FIND_FAILED is used when Benthos cannot finalize a "read"
	// operation on Profiles for MongoDB provider because of an unexpected error.
	CODE_BENTHOS_MONGODB_PROFILES_FIND_FAILED RemoraErrorCode = "04002"

	// CODE_BENTHOS_MONGODB_EVENTS_FIND_FAILED is used when Benthos cannot finalize a "read"
	// operation on Events for MongoDB provider because of an unexpected error.
	CODE_BENTHOS_MONGODB_EVENTS_FIND_FAILED RemoraErrorCode = "04003"

	// CODE_BENTHOS_MONGODB_EVENTS_FIND_FAILED is used when Benthos cannot finalize a "delete"
	// operation on Events for MongoDB provider because of an unexpected error.
	CODE_BENTHOS_MONGODB_EVENTS_DELETE_FAILED RemoraErrorCode = "04004"
)
