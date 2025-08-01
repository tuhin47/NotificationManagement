package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GeminiServiceImpl struct {
	*AIModelServiceImpl[models.GeminiModel, types.GeminiResponse]
	Repo        domain.GeminiModelRepository
	CurlService domain.CurlService
}

func NewGeminiService(repo domain.GeminiModelRepository, curlService domain.CurlService) domain.GeminiService {
	return &GeminiServiceImpl{
		AIModelServiceImpl: NewAIModelService[models.GeminiModel, types.GeminiResponse](repo),
		Repo:               repo,
		CurlService:        curlService,
	}
}

func (s *GeminiServiceImpl) MakeAIRequest(mod *models.AIModel, requestId uint) (*types.GeminiResponse, error) {
	curlResponse, err := s.CurlService.ExecuteCurlByRequestId(requestId)
	if err != nil {
		return nil, err
	}

	model, err := s.Repo.GetByID(context.Background(), mod.ID, nil)
	if err != nil {
		return nil, err
	}
	respBody, err := geminiCall(model, curlResponse)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	// Parse the response
	var geminiResp types.GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	return &geminiResp, nil
}

func geminiCall(model *models.GeminiModel, response *types.CurlResponse) ([]byte, error) {
	assistantContent, err := response.GetAssistantContent()
	if err != nil {
		return nil, err
	}

	geminiReq := types.GeminiRequest{
		Model: model.ModelName,
		Contents: []*types.GeminiMessage{
			{
				Role: "user", // Gemini typically starts with user role
				Parts: []types.GeminiPart{
					{Text: assistantContent},
				},
			},
			{
				Role: "model", // Gemini's assistant role is 'model'
				Parts: []types.GeminiPart{
					{Text: "Please check the current rate from the json.Is it greater than 125 ? Return Json Response "},
				},
			},
		},
	}

	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", model.GetBaseURL(), model.ModelName, model.APISecret)
	reqBody, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
