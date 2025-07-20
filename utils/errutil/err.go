package errutil

import (
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
	Message    string    `json:"message"`     // Human-readable error message
	Error      string    `json:"error"`       // Actual error message for debugging
	StatusCode int       `json:"status_code"` // HTTP status code
	Timestamp  time.Time `json:"timestamp"`   // Human-readable timestamp
	ErrorCode  string    `json:"error_code"`  // Internal error code for tracking
}

// ErrorCode represents different types of errors with their codes
type ErrorCode struct {
	Code    string
	Message string
	Status  int
}

// Predefined error codes
var (
	// ErrDatabaseConnection Database errors
	ErrDatabaseConnection = ErrorCode{
		Code:    "DB_CONNECTION_ERROR",
		Message: "Database connection failed",
		Status:  500,
	}
	ErrRecordNotFound = ErrorCode{
		Code:    "RECORD_NOT_FOUND",
		Message: "The requested record was not found",
		Status:  404,
	}
	ErrDatabaseQuery = ErrorCode{
		Code:    "DB_QUERY_ERROR",
		Message: "Database query failed",
		Status:  500,
	}

	// ErrInvalidInput Validation errors
	ErrInvalidInput = ErrorCode{
		Code:    "INVALID_INPUT",
		Message: "The provided input is invalid",
		Status:  400,
	}
	ErrInvalidEmail = ErrorCode{
		Code:    "INVALID_EMAIL",
		Message: "Invalid email format",
		Status:  400,
	}
	ErrInvalidPassword = ErrorCode{
		Code:    "INVALID_PASSWORD",
		Message: "Password does not meet requirements",
		Status:  400,
	}

	// ErrInvalidCredentials Authentication errors
	ErrInvalidCredentials = ErrorCode{
		Code:    "INVALID_CREDENTIALS",
		Message: "Invalid login credentials",
		Status:  401,
	}
	ErrUnauthorized = ErrorCode{
		Code:    "UNAUTHORIZED",
		Message: "Access denied",
		Status:  401,
	}
	ErrTokenExpired = ErrorCode{
		Code:    "TOKEN_EXPIRED",
		Message: "Authentication token has expired",
		Status:  401,
	}
	ErrInvalidToken = ErrorCode{
		Code:    "INVALID_TOKEN",
		Message: "Invalid authentication token",
		Status:  401,
	}

	// Business logic errors
	ErrUserAlreadyExists = ErrorCode{
		Code:    "USER_ALREADY_EXISTS",
		Message: "User already exists",
		Status:  409,
	}
	ErrUserNotFound = ErrorCode{
		Code:    "USER_NOT_FOUND",
		Message: "User not found",
		Status:  404,
	}
	ErrEmailAlreadyInUse = ErrorCode{
		Code:    "EMAIL_ALREADY_IN_USE",
		Message: "Email address is already in use",
		Status:  409,
	}

	// Server errors
	ErrInternalServer = ErrorCode{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
		Status:  500,
	}
	ErrServiceUnavailable = ErrorCode{
		Code:    "SERVICE_UNAVAILABLE",
		Message: "Service is temporarily unavailable",
		Status:  503,
	}

	// Notification errors
	ErrNotificationFailed = ErrorCode{
		Code:    "NOTIFICATION_FAILED",
		Message: "Failed to send notification",
		Status:  500,
	}
	ErrNotificationNotFound = ErrorCode{
		Code:    "NOTIFICATION_NOT_FOUND",
		Message: "Notification not found",
		Status:  404,
	}

	// File/Upload errors
	ErrFileUploadFailed = ErrorCode{
		Code:    "FILE_UPLOAD_FAILED",
		Message: "File upload failed",
		Status:  500,
	}
	ErrInvalidFileType = ErrorCode{
		Code:    "INVALID_FILE_TYPE",
		Message: "Invalid file type",
		Status:  400,
	}
	ErrFileTooLarge = ErrorCode{
		Code:    "FILE_TOO_LARGE",
		Message: "File size exceeds limit",
		Status:  400,
	}

	// Rate limiting
	ErrRateLimitExceeded = ErrorCode{
		Code:    "RATE_LIMIT_EXCEEDED",
		Message: "Too many requests",
		Status:  429,
	}

	// External service errors
	ErrExternalServiceError = ErrorCode{
		Code:    "EXTERNAL_SERVICE_ERROR",
		Message: "External service error",
		Status:  502,
	}
)

// CreateErrorResponse creates a new ErrorResponse with the given error code and actual error
func CreateErrorResponse(errCode ErrorCode, actualError error) ErrorResponse {
	return ErrorResponse{
		Message:    errCode.Message,
		Error:      actualError.Error(),
		StatusCode: errCode.Status,
		Timestamp:  time.Now(),
		ErrorCode:  errCode.Code,
	}
}

// CreateErrorResponseWithMessage creates a new ErrorResponse with custom message
func CreateErrorResponseWithMessage(errCode ErrorCode, actualError error, customMessage string) ErrorResponse {
	return ErrorResponse{
		Message:    customMessage,
		Error:      actualError.Error(),
		StatusCode: errCode.Status,
		Timestamp:  time.Now(),
		ErrorCode:  errCode.Code,
	}
}

// NewAppError creates a new AppError
func NewAppError(errCode ErrorCode, err error) error {
	return &AppError{
		Code:    errCode,
		Message: errCode.Message,
		Err:     err,
	}
}

// NewAppErrorWithMessage creates a new AppError with custom message
func NewAppErrorWithMessage(errCode ErrorCode, err error, message string) error {
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
			Timestamp:  time.Now(),
			ErrorCode:  appError.Code.Code,
		}
	}

	// If it's not an AppError, treat it as internal server error
	return CreateErrorResponse(ErrInternalServer, appErr)
}

// GetErrorCodeByCode returns an ErrorCode by its code string
func GetErrorCodeByCode(code string) (ErrorCode, bool) {
	errorCodes := map[string]ErrorCode{
		"DB_CONNECTION_ERROR":    ErrDatabaseConnection,
		"RECORD_NOT_FOUND":       ErrRecordNotFound,
		"DB_QUERY_ERROR":         ErrDatabaseQuery,
		"INVALID_INPUT":          ErrInvalidInput,
		"INVALID_EMAIL":          ErrInvalidEmail,
		"INVALID_PASSWORD":       ErrInvalidPassword,
		"INVALID_CREDENTIALS":    ErrInvalidCredentials,
		"UNAUTHORIZED":           ErrUnauthorized,
		"TOKEN_EXPIRED":          ErrTokenExpired,
		"INVALID_TOKEN":          ErrInvalidToken,
		"USER_ALREADY_EXISTS":    ErrUserAlreadyExists,
		"USER_NOT_FOUND":         ErrUserNotFound,
		"EMAIL_ALREADY_IN_USE":   ErrEmailAlreadyInUse,
		"INTERNAL_SERVER_ERROR":  ErrInternalServer,
		"SERVICE_UNAVAILABLE":    ErrServiceUnavailable,
		"NOTIFICATION_FAILED":    ErrNotificationFailed,
		"NOTIFICATION_NOT_FOUND": ErrNotificationNotFound,
		"FILE_UPLOAD_FAILED":     ErrFileUploadFailed,
		"INVALID_FILE_TYPE":      ErrInvalidFileType,
		"FILE_TOO_LARGE":         ErrFileTooLarge,
		"RATE_LIMIT_EXCEEDED":    ErrRateLimitExceeded,
		"EXTERNAL_SERVICE_ERROR": ErrExternalServiceError,
	}

	errCode, exists := errorCodes[code]
	return errCode, exists
}

// IsErrorCode checks if an error matches a specific error code
func IsErrorCode(err error, errCode ErrorCode) bool {
	// This is a simple implementation - you might want to enhance this
	// based on your specific error handling needs
	return err != nil && err.Error() == errCode.Message
}
