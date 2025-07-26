package errutil

import "net/http"

// Predefined error codes
var (
	// Database errors
	ErrDatabaseConnection = ErrorCode{Code: "DB_CONNECTION_ERROR", Message: "Database connection failed", Status: http.StatusInternalServerError}
	ErrDatabaseQuery      = ErrorCode{Code: "DB_QUERY_ERROR", Message: "Database query failed", Status: http.StatusInternalServerError}
	ErrRecordNotFound     = ErrorCode{Code: "RECORD_NOT_FOUND", Message: "The requested record was not found", Status: http.StatusNotFound}

	// Validation/Input errors
	ErrDuplicateEntry     = ErrorCode{Code: "DUPLICATE_ERROR", Message: "Duplicate Entry", Status: http.StatusBadRequest}
	ErrInvalidIdParam     = ErrorCode{Code: "INVALID_PARAM", Message: "Invalid Parameter", Status: http.StatusBadRequest}
	ErrInvalidRequestBody = ErrorCode{Code: "INVALID_BODY", Message: "Invalid Input", Status: http.StatusBadRequest}

	// Authentication/Authorization errors
	ErrInvalidCredentials = ErrorCode{Code: "INVALID_CREDENTIALS", Message: "Invalid login credentials", Status: http.StatusUnauthorized}
	ErrUnauthorized       = ErrorCode{Code: "UNAUTHORIZED", Message: "Access denied", Status: http.StatusUnauthorized}
	ErrTokenExpired       = ErrorCode{Code: "TOKEN_EXPIRED", Message: "Authentication token has expired", Status: http.StatusUnauthorized}
	ErrInvalidToken       = ErrorCode{Code: "INVALID_TOKEN", Message: "Invalid authentication token", Status: http.StatusUnauthorized}

	// Server/service errors
	ErrInternalServer     = ErrorCode{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error", Status: http.StatusInternalServerError}
	ErrServiceUnavailable = ErrorCode{Code: "SERVICE_UNAVAILABLE", Message: "Service is temporarily unavailable", Status: http.StatusServiceUnavailable}

	// Notification errors
	ErrNotificationFailed = ErrorCode{Code: "NOTIFICATION_FAILED", Message: "Failed to send notification", Status: http.StatusInternalServerError}

	// Rate limiting
	ErrRateLimitExceeded = ErrorCode{Code: "RATE_LIMIT_EXCEEDED", Message: "Too many requests", Status: http.StatusTooManyRequests}

	// External service errors
	ErrExternalServiceError = ErrorCode{Code: "EXTERNAL_SERVICE_ERROR", Message: "External service error", Status: http.StatusBadGateway}
)
