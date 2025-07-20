package throw

import (
	"NotificationManagement/middleware"
	"NotificationManagement/utils/errutil"
)

// ValidationError throws a validation error that will be caught by the global error handler
func ValidationError(field, message string) error {
	return middleware.NewValidationError(field, message)
}

// BusinessError throws a business error that will be caught by the global error handler
func BusinessError(code, message string) error {
	return middleware.NewBusinessError(code, message)
}

// AppError throws an AppError that will be caught by the global error handler
func AppError(errCode errutil.ErrorCode, err error) error {
	return errutil.NewAppError(errCode, err)
}

// AppErrorWithMessage throws an AppError with custom message
func AppErrorWithMessage(errCode errutil.ErrorCode, err error, message string) error {
	return errutil.NewAppErrorWithMessage(errCode, err, message)
}

// RecordNotFound throws a record not found error
func RecordNotFound(message string) error {
	if message == "" {
		message = "Record not found"
	}
	return errutil.NewAppError(errutil.ErrRecordNotFound, errutil.NewError(message))
}

// InvalidInput throws an invalid input error
func InvalidInput(message string) error {
	if message == "" {
		message = "Invalid input provided"
	}
	return errutil.NewAppError(errutil.ErrInvalidInput, errutil.NewError(message))
}

// Unauthorized throws an unauthorized error
func Unauthorized(message string) error {
	if message == "" {
		message = "Access denied"
	}
	return errutil.NewAppError(errutil.ErrUnauthorized, errutil.NewError(message))
}

// InternalServerError throws an internal server error
func InternalServerError(message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return errutil.NewAppError(errutil.ErrInternalServer, errutil.NewError(message))
}
