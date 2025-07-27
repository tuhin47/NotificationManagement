package controllers

import (
	"NotificationManagement/domain"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"NotificationManagement/utils/errutil"
	"NotificationManagement/utils/throw"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AIControllerImpl struct {
	OllamaService   domain.OllamaService
	CurlService     domain.CurlService
	DeepseekService domain.DeepseekModelService
}

func NewAIController(ollamaService domain.OllamaService, curlService domain.CurlService, deepseekService domain.DeepseekModelService) domain.AIController {
	return &AIControllerImpl{
		OllamaService:   ollamaService,
		CurlService:     curlService,
		DeepseekService: deepseekService,
	}
}

func (ac *AIControllerImpl) MakeAIRequestHandler(c echo.Context) error {
	var req types.MakeAIRequestPayload
	if err := utils.BindAndValidate(c, &req); err != nil {
		return err
	}

	// Get the CurlRequest by ID
	curlRequest, err := ac.CurlService.GetCurlRequestByID(req.CurlRequestID)
	if err != nil {
		return throw.AppError(errutil.ErrRecordNotFound, err)
	}

	// Convert CurlRequest model to types.CurlRequest
	curlReq := types.CurlRequest{
		URL:     curlRequest.URL,
		Method:  curlRequest.Method,
		Body:    curlRequest.Body,
		RawCurl: curlRequest.RawCurl,
	}

	// Parse headers from string to map
	if curlRequest.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(curlRequest.Headers), &headers); err != nil {
			return throw.AppError(errutil.ErrInvalidRequestBody, err)
		}
		curlReq.Headers = headers
	}

	// Execute the curl request to get CurlResponse
	curlResponse, err := ac.CurlService.ExecuteCurl(curlReq)
	if err != nil {
		return err
	}

	// Get the DeepseekModel by ID
	deepseekModel, err := ac.DeepseekService.GetDeepseekModelByID(req.ModelID)
	if err != nil {
		return throw.AppError(errutil.ErrRecordNotFound, err)
	}

	// Make the AI request using OllamaService
	ollamaResponse, err := ac.OllamaService.MakeAIRequest(*deepseekModel, curlResponse)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ollamaResponse)
}
