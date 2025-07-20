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
