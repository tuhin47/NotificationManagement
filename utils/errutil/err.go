package errutil

import (
	"NotificationManagement/logger"
	"errors"
	"time"
)

type ErrorCode struct {
	Code    string
	Message string
	Status  int
}

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

type ErrorResponse struct {
	Message    string    `json:"message"`
	Error      string    `json:"error"`
	StatusCode int       `json:"status_code"`
	Timestamp  time.Time `json:"timestamp"`
	ErrorCode  string    `json:"error_code"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(errCode ErrorCode, err error) error {
	return NewAppErrorWithMessage(errCode, err, errCode.Message)
}

func NewAppErrorWithMessage(errCode ErrorCode, err error, message string) error {
	var target *AppError
	if errors.As(err, &target) {
		logger.DPanic("Duplicate Throws", err)
		return err
	}
	if err == nil {
		err = ErrUndefine
	}
	logger.DPanic("Error occurred", err)
	return &AppError{
		Code:    errCode,
		Message: message,
		Err:     err,
	}
}

func CreateErrorResponse(errCode ErrorCode, actualError error) ErrorResponse {
	return ErrorResponse{
		Message:    errCode.Message,
		Error:      actualError.Error(),
		StatusCode: errCode.Status,
		Timestamp:  GetCurrentTime(),
		ErrorCode:  errCode.Code,
	}
}

func CreateErrorResponseWithMessage(errCode ErrorCode, actualError error, customMessage string) ErrorResponse {
	return ErrorResponse{
		Message:    customMessage,
		Error:      actualError.Error(),
		StatusCode: errCode.Status,
		Timestamp:  GetCurrentTime(),
		ErrorCode:  errCode.Code,
	}
}

func AppErrorToErrorResponse(appErr error) ErrorResponse {
	var appError *AppError
	if errors.As(appErr, &appError) {
		return ErrorResponse{
			Message:    appError.Message,
			Error:      appError.Err.Error(),
			StatusCode: appError.Code.Status,
			Timestamp:  GetCurrentTime(),
			ErrorCode:  appError.Code.Code,
		}
	}

	return CreateErrorResponse(ErrInternalServer, appErr)
}

func GetCurrentTime() time.Time {
	return time.Now()
}
