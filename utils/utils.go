package utils

import (
	"NotificationManagement/utils/errutil"
	"strconv"

	"github.com/labstack/echo/v4"
)

// BindAndValidate binds the request body to the target struct and returns a standardized error if binding fails.
func BindAndValidate(c echo.Context, target interface{}) error {
	if err := c.Bind(target); err != nil {
		return errutil.NewAppError(errutil.ErrInvalidRequestBody, err)
	}
	return nil
}

// ParseIDFromContext parses the "id" parameter from the echo.Context and returns it as a uint.
func ParseIDFromContext(c echo.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}
	return uint(id), nil
}

// ParseLimitAndOffset parses "limit" and "offset" query parameters from the echo.Context.
// It returns the parsed limit and offset, with default values of 10 and 0 respectively.
func ParseLimitAndOffset(c echo.Context) (limit, offset int) {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit = 10 // default limit
	offset = 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	return limit, offset
}
