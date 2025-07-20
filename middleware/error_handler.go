package middleware

import (
	"NotificationManagement/utils/errutil"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorHandler is a global error handler middleware that mimics Java's @ControllerAdvice
func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Execute the handler
			err := next(c)
			if err == nil {
				return nil
			}

			// Handle different types of errors
			return handleError(c, err)
		}
	}
}

// handleError processes different types of errors and returns appropriate HTTP responses
func handleError(c echo.Context, err error) error {
	// Check if it's an AppError
	if appError, ok := err.(*errutil.AppError); ok {
		errResp := errutil.AppErrorToErrorResponse(appError)
		return c.JSON(errResp.StatusCode, errResp)
	}

	// Check if it's a ValidationError
	if validationError, ok := err.(*ValidationError); ok {
		errResp := errutil.ErrorResponse{
			Message:    validationError.Message,
			Error:      err.Error(),
			StatusCode: 400,
			Timestamp:  errutil.GetCurrentTime(),
			ErrorCode:  "VALIDATION_ERROR",
		}
		return c.JSON(400, errResp)
	}

	// Check if it's a BusinessError
	if businessError, ok := err.(*BusinessError); ok {
		errResp := errutil.ErrorResponse{
			Message:    businessError.Message,
			Error:      err.Error(),
			StatusCode: 400,
			Timestamp:  errutil.GetCurrentTime(),
			ErrorCode:  businessError.Code,
		}
		return c.JSON(400, errResp)
	}

	// Check if it's an HTTPError (Echo's built-in error)
	if httpError, ok := err.(*echo.HTTPError); ok {
		errResp := errutil.ErrorResponse{
			Message:    getErrorMessage(err),
			Error:      err.Error(),
			StatusCode: httpError.Code,
			Timestamp:  errutil.GetCurrentTime(),
			ErrorCode:  "HTTP_ERROR",
		}
		return c.JSON(httpError.Code, errResp)
	}

	// Check for specific error types and map them to appropriate responses
	if errutil.IsRecordNotFound(err) {
		errResp := errutil.CreateErrorResponse(errutil.ErrRecordNotFound, err)
		return c.JSON(errResp.StatusCode, errResp)
	}

	if errutil.IsInvalidInput(err) {
		errResp := errutil.CreateErrorResponse(errutil.ErrInvalidInput, err)
		return c.JSON(errResp.StatusCode, errResp)
	}

	// Default to internal server error
	errResp := errutil.CreateErrorResponse(errutil.ErrInternalServer, err)
	return c.JSON(http.StatusInternalServerError, errResp)
}

func getErrorMessage(err error) string {
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		return "Bad Request"
	}

	return err.Error()

}

// ValidationError represents validation errors that can be thrown from controllers
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// BusinessError represents business logic errors that can be thrown from services
type BusinessError struct {
	Code    string
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

// NewBusinessError creates a new business error
func NewBusinessError(code, message string) error {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}
