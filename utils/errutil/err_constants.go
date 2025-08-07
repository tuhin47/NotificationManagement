package errutil

import (
	"errors"
	"net/http"
)

// Auth specific errors
var (
	ErrNoAuthHeader          = errors.New("no authorization header provided")
	ErrInvalidHeaderFormat   = errors.New("invalid authorization header format")
	ErrInvalidTokenSignature = errors.New("invalid token signing method")
	ErrInvalidTokenValue     = errors.New("invalid token value")
	ErrNoCertificateFound    = errors.New("no valid certificate found")
	ErrNoRoleInformation     = errors.New("no role information in context")
	ErrInsufficientRoles     = errors.New("user does not have required roles")
	ErrInvalidFeature        = errors.New("invalid feature")
)

// Predefined error codes
var (
	// Database errors
	ErrDatabaseQuery  = ErrorCode{Code: "DB_QUERY_ERROR", Message: "Database query failed", Status: http.StatusInternalServerError}
	ErrRecordNotFound = ErrorCode{Code: "RECORD_NOT_FOUND", Message: "The requested record was not found", Status: http.StatusNotFound}

	// Validation/Input errors
	ErrDuplicateEntry      = ErrorCode{Code: "DUPLICATE_ERROR", Message: "Duplicate Entry", Status: http.StatusBadRequest}
	ErrInvalidIdParam      = ErrorCode{Code: "INVALID_PARAM", Message: "Invalid Parameter", Status: http.StatusBadRequest}
	ErrInvalidRequestBody  = ErrorCode{Code: "INVALID_BODY", Message: "Invalid Input", Status: http.StatusBadRequest}
	ErrFeatureNotAvailable = ErrorCode{Code: "INVALID_FEATURE", Message: "Feature not available", Status: http.StatusNotImplemented}

	// Authentication/Authorization errors
	ErrInvalidToken           = ErrorCode{Code: "INVALID_TOKEN", Message: "Invalid authentication token", Status: http.StatusUnauthorized}
	ErrMissingAuthHeader      = ErrorCode{Code: "MISSING_AUTH_HEADER", Message: "Missing Authorization header", Status: http.StatusUnauthorized}
	ErrInvalidAuthFormat      = ErrorCode{Code: "INVALID_AUTH_FORMAT", Message: "Invalid Authorization header format", Status: http.StatusUnauthorized}
	ErrInvalidSigningMethod   = ErrorCode{Code: "INVALID_SIGNING_METHOD", Message: "Unexpected token signing method", Status: http.StatusUnauthorized}
	ErrCertificateRetrieval   = ErrorCode{Code: "CERT_RETRIEVAL_ERROR", Message: "Failed to retrieve certificate", Status: http.StatusUnauthorized}
	ErrNoCertificates         = ErrorCode{Code: "NO_CERTIFICATES", Message: "No certificates found", Status: http.StatusUnauthorized}
	ErrNoRoleInfo             = ErrorCode{Code: "NO_ROLE_INFO", Message: "No role information available", Status: http.StatusForbidden}
	ErrInsufficientPrivileges = ErrorCode{Code: "INSUFFICIENT_PRIVILEGES", Message: "Access denied: Insufficient privileges", Status: http.StatusForbidden}
	ErrUserRegistrationFailed = ErrorCode{Code: "USER_REGISTRATION_FAILED", Message: "Failed to register or update user", Status: http.StatusInternalServerError}

	// Server/service errors
	ErrInternalServer = ErrorCode{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error", Status: http.StatusInternalServerError}

	// Task/Queue errors
	ErrTaskEnqueueFailed        = ErrorCode{Code: "TASK_ENQUEUE_FAILED", Message: "Failed to enqueue task", Status: http.StatusInternalServerError}
	ErrTaskDeletionFailed       = ErrorCode{Code: "TASK_DELETION_FAILED", Message: "Failed to delete task", Status: http.StatusInternalServerError}
	ErrTaskInfoRetrievalFailed  = ErrorCode{Code: "TASK_INFO_RETRIEVAL_FAILED", Message: "Failed to retrieve task info", Status: http.StatusInternalServerError}
	ErrTaskMarshalPayloadFailed = ErrorCode{Code: "TASK_MARSHAL_PAYLOAD_FAILED", Message: "Failed to marshal task payload", Status: http.StatusInternalServerError}

	// Rate limiting

	// External service errors
	ErrExternalServiceError       = ErrorCode{Code: "EXTERNAL_SERVICE_ERROR", Message: "External service error", Status: http.StatusBadGateway}
	ErrCurlCommandExecutionFailed = ErrorCode{Code: "CURL_COMMAND_EXECUTION_FAILED", Message: "Failed to execute curl command", Status: http.StatusInternalServerError}
	ErrCurlParseError             = ErrorCode{Code: "CURL_PARSE_ERROR", Message: "Failed to parse curl command", Status: http.StatusBadRequest}

	// AI Model specific errors
	ErrUnsupportedAIModelType = ErrorCode{Code: "UNSUPPORTED_AI_MODEL_TYPE", Message: "Unsupported AI model type", Status: http.StatusBadRequest}
	ErrAIMarshalRequestFailed = ErrorCode{Code: "AI_MARSHAL_REQUEST_FAILED", Message: "Failed to marshal AI request", Status: http.StatusInternalServerError}
	ErrAICreateRequestFailed  = ErrorCode{Code: "AI_CREATE_REQUEST_FAILED", Message: "Failed to create AI HTTP request", Status: http.StatusInternalServerError}
	ErrAIPullModelFailed      = ErrorCode{Code: "AI_PULL_MODEL_FAILED", Message: "Failed to pull AI model", Status: http.StatusInternalServerError}

	// Curl Response specific errors
	ErrEmptyResponse                 = ErrorCode{Code: "EMPTY_RESPONSE", Message: "Empty response", Status: http.StatusInternalServerError}
	ErrCurlMarshalResponseBodyFailed = ErrorCode{Code: "MARSHAL_RESPONSE_BODY_FAILED", Message: "Failed to marshal request body", Status: http.StatusInternalServerError}
	ErrCurlInvalidResponseBodyType   = ErrorCode{Code: "INVALID_RESPONSE_BODY_TYPE", Message: "Invalid request body type", Status: http.StatusInternalServerError}
	ErrCurlCreateTempFileFailed      = ErrorCode{Code: "CREATE_TEMP_FILE_FAILED", Message: "Failed to create temporary file for request", Status: http.StatusInternalServerError}
	ErrCurlWriteTempFileFailed       = ErrorCode{Code: "WRITE_TEMP_FILE_FAILED", Message: "Failed to write to temporary file for request", Status: http.StatusInternalServerError}
	ErrCurlUnsupportedResponseType   = ErrorCode{Code: "UNSUPPORTED_RESPONSE_TYPE", Message: "Unsupported request type", Status: http.StatusBadRequest}

	// GORM utility errors
	ErrGormInvalidSlicePointer = ErrorCode{Code: "GORM_INVALID_SLICE_POINTER", Message: "New items must be a pointer to a slice", Status: http.StatusInternalServerError}
	ErrGormAssociationNotFound = ErrorCode{Code: "GORM_ASSOCIATION_NOT_FOUND", Message: "GORM association not found", Status: http.StatusInternalServerError}
	ErrGormCreateFailed        = ErrorCode{Code: "GORM_CREATE_FAILED", Message: "Failed to create record in GORM", Status: http.StatusInternalServerError}
	ErrGormUpdateFailed        = ErrorCode{Code: "GORM_UPDATE_FAILED", Message: "Failed to update record in GORM", Status: http.StatusInternalServerError}
	ErrGormDeleteFailed        = ErrorCode{Code: "GORM_DELETE_FAILED", Message: "Failed to delete record in GORM", Status: http.StatusInternalServerError}
	ErrGormFindFailed          = ErrorCode{Code: "GORM_FIND_FAILED", Message: "Failed to find record in GORM", Status: http.StatusInternalServerError}
)
