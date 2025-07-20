package errutil

import (
	"NotificationManagement/logger"
	"errors"
	"net/http"
	"time"
)

// AppError represents a custom error with additional context
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// ErrorResponse represents a standardized error response structure
type ErrorResponse struct {
	Message    string    `json:"message"`
	Error      string    `json:"error"`
	StatusCode int       `json:"status_code"`
	Timestamp  time.Time `json:"timestamp"`
	ErrorCode  string    `json:"error_code"`
}

// ErrorCode represents different types of errors with their codes
type ErrorCode struct {
	Code    string
	Message string
	Status  int
}

// Predefined error codes
var (
	// Database errors
	ErrDatabaseConnection = ErrorCode{Code: "DB_CONNECTION_ERROR", Message: "Database connection failed", Status: http.StatusInternalServerError}
	ErrDatabaseQuery      = ErrorCode{Code: "DB_QUERY_ERROR", Message: "Database query failed", Status: http.StatusInternalServerError}
	ErrRecordNotFound     = ErrorCode{Code: "RECORD_NOT_FOUND", Message: "The requested record was not found", Status: http.StatusNotFound}

	// Validation/Input errors
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

// CreateErrorResponse creates a new ErrorResponse with the given error code and actual error
func CreateErrorResponse(errCode ErrorCode, actualError error) ErrorResponse {
	return ErrorResponse{
		Message:    errCode.Message,
		Error:      actualError.Error(),
		StatusCode: errCode.Status,
		Timestamp:  GetCurrentTime(),
		ErrorCode:  errCode.Code,
	}
}

// CreateErrorResponseWithMessage creates a new ErrorResponse with custom message
func CreateErrorResponseWithMessage(errCode ErrorCode, actualError error, customMessage string) ErrorResponse {
	return ErrorResponse{
		Message:    customMessage,
		Error:      actualError.Error(),
		StatusCode: errCode.Status,
		Timestamp:  GetCurrentTime(),
		ErrorCode:  errCode.Code,
	}
}

// NewAppError creates a new AppError
func NewAppError(errCode ErrorCode, err error) error {
	var target *AppError
	if errors.As(err, &target) {
		return err
	}
	logger.DPanic("Error occurred", err)
	return &AppError{
		Code:    errCode,
		Message: errCode.Message,
		Err:     err,
	}
}

// NewAppErrorWithMessage creates a new AppError with custom message
func NewAppErrorWithMessage(errCode ErrorCode, err error, message string) error {
	logger.Error("Error occurred", err)
	return &AppError{
		Code:    errCode,
		Message: message,
		Err:     err,
	}
}

// AppErrorToErrorResponse converts an AppError to ErrorResponse
func AppErrorToErrorResponse(appErr error) ErrorResponse {
	if appError, ok := appErr.(*AppError); ok {
		return ErrorResponse{
			Message:    appError.Message,
			Error:      appError.Err.Error(),
			StatusCode: appError.Code.Status,
			Timestamp:  GetCurrentTime(),
			ErrorCode:  appError.Code.Code,
		}
	}

	// If it's not an AppError, treat it as internal server error
	return CreateErrorResponse(ErrInternalServer, appErr)
}

// GetCurrentTime returns the current time for error responses
func GetCurrentTime() time.Time {
	return time.Now()
}

// NewError creates a new error with the given message
func NewError(message string) error {
	return errors.New(message)
}
