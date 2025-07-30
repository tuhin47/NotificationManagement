package utils

import (
	"NotificationManagement/utils/errutil"
	"github.com/labstack/echo/v4"
)

// BindAndValidate binds the request body to the target struct and returns a standardized error if binding fails.
func BindAndValidate(c echo.Context, target interface{}) error {
	if err := c.Bind(target); err != nil {
		return errutil.NewAppError(errutil.ErrInvalidRequestBody, err)
	}
	return nil
}
