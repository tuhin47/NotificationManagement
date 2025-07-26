package middleware

import (
	"NotificationManagement/types"
	"NotificationManagement/utils"
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
			err := next(c)
			if err == nil {
				return nil
			}
			return handleError(c, err)
		}
	}
}

// handleError processes different types of errors and returns appropriate HTTP responses
func handleError(c echo.Context, err error) error {

	// Check if it's an AppError
	var appError *errutil.AppError
	if errors.As(err, &appError) {
		errResp := errutil.AppErrorToErrorResponse(appError)
		return c.JSON(errResp.StatusCode, errResp)
	}

	// Check if it's a BusinessError
	var businessError *BusinessError
	if errors.As(err, &businessError) {
		errResp := types.ErrorResponse{
			Message:    businessError.Message,
			Error:      err.Error(),
			StatusCode: 400,
			Timestamp:  utils.GetCurrentTime(),
			ErrorCode:  businessError.Code,
		}
		return c.JSON(400, errResp)
	}

	// Check if it's an HTTPError (Echo's built-in error)
	if httpError, ok := err.(*echo.HTTPError); ok {
		errResp := types.ErrorResponse{
			Message:    getErrorMessage(err),
			Error:      err.Error(),
			StatusCode: httpError.Code,
			Timestamp:  utils.GetCurrentTime(),
			ErrorCode:  "HTTP_ERROR",
		}
		return c.JSON(httpError.Code, errResp)
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
