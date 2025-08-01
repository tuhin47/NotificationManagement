package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AIRequestControllerImpl struct {
	ServiceFact domain.AIServiceManager
}

func NewAIRequestController(s domain.AIServiceManager) domain.AIRequestController {
	return &AIRequestControllerImpl{
		ServiceFact: s,
	}
}

func (ac *AIRequestControllerImpl) MakeAIRequestHandler(c echo.Context) error {
	var req types.MakeAIRequestPayload
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}
	aiResponse, err := ac.ServiceFact.ProcessAIRequest(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, aiResponse)
}
