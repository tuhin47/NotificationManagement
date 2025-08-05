package errutil

import (
	"errors"
	"net/http"
)

// Auth specific errors
var (
	ErrNoAuthHeader                = errors.New("no authorization header provided")
	ErrInvalidHeaderFormat         = errors.New("invalid authorization header format")
	ErrInvalidTokenSignature       = errors.New("invalid token signing method")
	ErrInvalidTokenValue           = errors.New("invalid token value")
	ErrNoCertificateFound          = errors.New("no valid certificate found")
	ErrNoRoleInformation           = errors.New("no role information in context")
	ErrInsufficientRoles           = errors.New("user does not have required roles")
	ErrServiceNotAvailable         = errors.New("service not available")
	ErrInvalidFeature              = errors.New("invalid feature")
	ErrFailedToCastDeepseekModel   = errors.New("failed to cast model to DeepseekModel")
	ErrFailedToCastGeminiModel     = errors.New("failed to cast model to GeminiModel")
	ErrUnsupportedAIModelTypeMsg   = errors.New("unsupported AI model type")
	ErrAIModelNotFoundMsg          = errors.New("AI Model not found")
	ErrUnknownAIModelTypeMsg       = errors.New("unknown AI model type")
	ErrAIModelNotFoundOrDeletedMsg = errors.New("AI Model not found or could not be deleted")
)

// Predefined error codes
var (
	// Database errors
	ErrDatabaseConnection = ErrorCode{Code: "DB_CONNECTION_ERROR", Message: "Database connection failed", Status: http.StatusInternalServerError}
	ErrDatabaseQuery      = ErrorCode{Code: "DB_QUERY_ERROR", Message: "Database query failed", Status: http.StatusInternalServerError}
	ErrRecordNotFound     = ErrorCode{Code: "RECORD_NOT_FOUND", Message: "The requested record was not found", Status: http.StatusNotFound}

	// Validation/Input errors
	ErrDuplicateEntry      = ErrorCode{Code: "DUPLICATE_ERROR", Message: "Duplicate Entry", Status: http.StatusBadRequest}
	ErrInvalidIdParam      = ErrorCode{Code: "INVALID_PARAM", Message: "Invalid Parameter", Status: http.StatusBadRequest}
	ErrInvalidRequestBody  = ErrorCode{Code: "INVALID_BODY", Message: "Invalid Input", Status: http.StatusBadRequest}
	ErrFeatureNotAvailable = ErrorCode{Code: "INVALID_FEATURE", Message: "Feature not available", Status: http.StatusNotImplemented}

	// Authentication/Authorization errors
	ErrInvalidCredentials     = ErrorCode{Code: "INVALID_CREDENTIALS", Message: "Invalid login credentials", Status: http.StatusUnauthorized}
	ErrUnauthorized           = ErrorCode{Code: "UNAUTHORIZED", Message: "Access denied", Status: http.StatusUnauthorized}
	ErrTokenExpired           = ErrorCode{Code: "TOKEN_EXPIRED", Message: "Authentication token has expired", Status: http.StatusUnauthorized}
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
	ErrInternalServer     = ErrorCode{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error", Status: http.StatusInternalServerError}
	ErrServiceUnavailable = ErrorCode{Code: "SERVICE_UNAVAILABLE", Message: "ServiceFact is temporarily unavailable", Status: http.StatusServiceUnavailable}

	// Notification errors
	ErrNotificationFailed = ErrorCode{Code: "NOTIFICATION_FAILED", Message: "Failed to send notification", Status: http.StatusInternalServerError}

	// Rate limiting
	ErrRateLimitExceeded = ErrorCode{Code: "RATE_LIMIT_EXCEEDED", Message: "Too many requests", Status: http.StatusTooManyRequests}

	// External service errors
	ErrExternalServiceError = ErrorCode{Code: "EXTERNAL_SERVICE_ERROR", Message: "External service error", Status: http.StatusBadGateway}

	// AI Model specific errors
	ErrFailedToCastModel      = ErrorCode{Code: "AI_MODEL_CAST_ERROR", Message: "Failed to cast AI model", Status: http.StatusInternalServerError}
	ErrUnsupportedAIModelType = ErrorCode{Code: "UNSUPPORTED_AI_MODEL_TYPE", Message: "Unsupported AI model type", Status: http.StatusBadRequest}
)
