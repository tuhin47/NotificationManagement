package errutil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}
			return HandleError(c, err)
		}
	}
}

func HandleError(c echo.Context, err error) error {

	var appError *AppError
	if errors.As(err, &appError) {
		errResp := AppErrorToErrorResponse(appError)
		return c.JSON(errResp.StatusCode, errResp)
	}

	var businessError *BusinessError
	if errors.As(err, &businessError) {
		errResp := ErrorResponse{
			Message:    businessError.Message,
			Error:      err.Error(),
			StatusCode: 400,
			Timestamp:  GetCurrentTime(),
			ErrorCode:  businessError.Code,
		}
		return c.JSON(400, errResp)
	}

	if httpError, ok := err.(*echo.HTTPError); ok {
		errResp := ErrorResponse{
			Message:    getErrorMessage(err),
			Error:      err.Error(),
			StatusCode: httpError.Code,
			Timestamp:  GetCurrentTime(),
			ErrorCode:  "HTTP_ERROR",
		}
		return c.JSON(httpError.Code, errResp)
	}

	errResp := CreateErrorResponse(ErrInternalServer, err)
	return c.JSON(http.StatusInternalServerError, errResp)
}

func getErrorMessage(err error) string {
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		return "Bad Request"
	}

	return err.Error()

}

type BusinessError struct {
	Code    string
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(code, message string) error {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}
