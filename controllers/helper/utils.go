package helper

import (
	"NotificationManagement/middleware"
	"NotificationManagement/utils/errutil"
	"github.com/labstack/echo/v4"
	"strconv"
)

func BindAndValidate(c echo.Context, target interface{}) error {
	if err := c.Bind(target); err != nil {
		return errutil.NewAppError(errutil.ErrInvalidRequestBody, err)
	}
	return nil
}

func ParseIDFromContext(c echo.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errutil.NewAppError(errutil.ErrInvalidIdParam, err)
	}
	return uint(id), nil
}

func ParseLimitAndOffset(c echo.Context) (limit, offset int) {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit = 10
	offset = 0

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

func GetUserId(c echo.Context) uint {
	ccx, _ := c.(*middleware.CustomContext)
	return ccx.UserID
}
