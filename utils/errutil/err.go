package errutil

import (
	"NotificationManagement/logger"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"errors"
)

// ErrorCode represents different types of errors with their codes
type ErrorCode struct {
	Code    string
	Message string
	Status  int
}

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

// NewAppError creates a new AppError
func NewAppError(errCode ErrorCode, err error) error {
	return NewAppErrorWithMessage(errCode, err, errCode.Message)
}

// NewAppErrorWithMessage creates a new AppError with custom message
func NewAppErrorWithMessage(errCode ErrorCode, err error, message string) error {
	var target *AppError
	if errors.As(err, &target) {
		logger.DPanic("Duplicate Throws", err)
		return err
	}
	logger.DPanic("Error occurred", err)
	return &AppError{
		Code:    errCode,
		Message: message,
		Err:     err,
	}
}

// CreateErrorResponse creates a new ErrorResponse with the given error code and actual error
func CreateErrorResponse(errCode ErrorCode, actualError error) types.ErrorResponse {
	return types.ErrorResponse{
		Message:    errCode.Message,
		Error:      actualError.Error(),
		StatusCode: errCode.Status,
		Timestamp:  utils.GetCurrentTime(),
		ErrorCode:  errCode.Code,
	}
}

// CreateErrorResponseWithMessage creates a new ErrorResponse with custom message
func CreateErrorResponseWithMessage(errCode ErrorCode, actualError error, customMessage string) types.ErrorResponse {
	return types.ErrorResponse{
		Message:    customMessage,
		Error:      actualError.Error(),
		StatusCode: errCode.Status,
		Timestamp:  utils.GetCurrentTime(),
		ErrorCode:  errCode.Code,
	}
}

// AppErrorToErrorResponse converts an AppError to ErrorResponse
func AppErrorToErrorResponse(appErr error) types.ErrorResponse {
	var appError *AppError
	if errors.As(appErr, &appError) {
		return types.ErrorResponse{
			Message:    appError.Message,
			Error:      appError.Err.Error(),
			StatusCode: appError.Code.Status,
			Timestamp:  utils.GetCurrentTime(),
			ErrorCode:  appError.Code.Code,
		}
	}

	// If it's not an AppError, treat it as internal server error
	return CreateErrorResponse(ErrInternalServer, appErr)
}
