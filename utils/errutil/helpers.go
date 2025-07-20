package errutil

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

// Common error checking functions
var (
	ErrRecordNotFoundError = errors.New("record not found")
	ErrInvalidInputError   = errors.New("invalid input")
)

// IsRecordNotFound checks if the error is a "record not found" error
func IsRecordNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, ErrRecordNotFoundError) ||
		(err != nil && err.Error() == "record not found")
}

// IsInvalidInput checks if the error is an invalid input error
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInputError) ||
		(err != nil && err.Error() == "invalid input")
}

// HandleDatabaseError converts database errors to appropriate ErrorResponse
func HandleDatabaseError(err error) ErrorResponse {
	if IsRecordNotFound(err) {
		return CreateErrorResponse(ErrRecordNotFound, err)
	}
	return CreateErrorResponse(ErrDatabaseQuery, err)
}

// HandleValidationError converts validation errors to appropriate ErrorResponse
func HandleValidationError(err error) ErrorResponse {
	return CreateErrorResponse(ErrInvalidInput, err)
}

// HandleAuthenticationError converts authentication errors to appropriate ErrorResponse
func HandleAuthenticationError(err error) ErrorResponse {
	return CreateErrorResponse(ErrInvalidCredentials, err)
}

// HandleAuthorizationError converts authorization errors to appropriate ErrorResponse
func HandleAuthorizationError(err error) ErrorResponse {
	return CreateErrorResponse(ErrUnauthorized, err)
}

// HandleInternalError converts internal errors to appropriate ErrorResponse
func HandleInternalError(err error) ErrorResponse {
	return CreateErrorResponse(ErrInternalServer, err)
}

// HandleAppError converts an AppError to ErrorResponse
func HandleAppError(err error) ErrorResponse {
	return AppErrorToErrorResponse(err)
}

// HandleServiceError handles errors from services and converts them to ErrorResponse
func HandleServiceError(err error) ErrorResponse {
	if err == nil {
		return CreateErrorResponse(ErrInternalServer, errors.New("unknown error"))
	}

	// Check if it's an AppError
	if appError, ok := err.(*AppError); ok {
		return AppErrorToErrorResponse(appError)
	}

	// Check for specific error types
	if IsRecordNotFound(err) {
		return CreateErrorResponse(ErrRecordNotFound, err)
	}

	if IsInvalidInput(err) {
		return CreateErrorResponse(ErrInvalidInput, err)
	}

	// Default to internal server error
	return CreateErrorResponse(ErrInternalServer, err)
}

// WriteErrorResponse writes an ErrorResponse to HTTP response writer
func WriteErrorResponse(w http.ResponseWriter, errResp ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResp.StatusCode)

	// Marshal the error response to JSON
	jsonData, err := json.Marshal(errResp)
	if err != nil {
		// Fallback to simple error response if marshaling fails
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal error response"}`))
		return
	}

	w.Write(jsonData)
}

// CreateErrorResponseFromError creates an ErrorResponse from a standard error
func CreateErrorResponseFromError(err error) ErrorResponse {
	if err == nil {
		return CreateErrorResponse(ErrInternalServer, errors.New("unknown error"))
	}

	// Check for specific error types
	if IsRecordNotFound(err) {
		return CreateErrorResponse(ErrRecordNotFound, err)
	}

	if IsInvalidInput(err) {
		return CreateErrorResponse(ErrInvalidInput, err)
	}

	// Default to internal server error
	return CreateErrorResponse(ErrInternalServer, err)
}

// ValidateRequiredFields checks if required fields are present
func ValidateRequiredFields(fields map[string]interface{}) error {
	for fieldName, value := range fields {
		if value == nil || value == "" {
			return errors.New(fieldName + " is required")
		}
	}
	return nil
}

// ValidateEmailFormat validates email format (basic validation)
func ValidateEmailFormat(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// Basic email validation - you might want to use a proper email validation library
	if len(email) < 5 || !contains(email, "@") {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Add more password validation rules as needed
	// For example: check for uppercase, lowercase, numbers, special characters

	return nil
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
