package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/repositories"
	"NotificationManagement/types"
	"NotificationManagement/types/ollama"
	"NotificationManagement/utils/errutil"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DeepseekServiceImpl struct {
	domain.CommonService[models.DeepseekModel]
	CurlService domain.CurlService
}

func NewDeepseekModelService(repo domain.DeepseekModelRepository, curl domain.CurlService) domain.DeepseekService {
	service := &DeepseekServiceImpl{
		CurlService: curl,
	}
	service.CommonService = NewCommonService(repo, service)
	return service
}

func (s *DeepseekServiceImpl) GetContext() context.Context {
	background := context.Background()
	f := []repositories.Filter{
		{Field: "type", Op: "=", Value: "deepseek"},
	}
	return context.WithValue(background, repositories.ContextStruct{}, &repositories.ContextStruct{Filter: &f})
}

func (s *DeepseekServiceImpl) MakeAIRequest(mod *models.AIModel, requestId uint) (interface{}, error) {
	curl, err := s.CurlService.GetModelByID(requestId)
	if err != nil {
		return nil, err
	}
	curlResponse, err := s.CurlService.ExecuteCurl(curl)
	if err != nil {
		return nil, err
	}
	model, err := s.GetModelByID(mod.ID)
	if err != nil {
		return nil, err
	}
	respBody, err := deepseekCall(model, curlResponse, curl)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	// Parse the response
	var ollamaResp ollama.OllamaResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	a := any(ollamaResp)
	return &a, nil
}

func (s *DeepseekServiceImpl) PullModel(model *models.DeepseekModel) error {
	// Implementation for pulling/downloading the model
	payload := map[string]interface{}{
		"name": model.ModelName,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal pull request: %w", err)
	}

	url := fmt.Sprintf("%s/api/pull", model.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to pull model, status code: %d", resp.StatusCode)
	}

	return nil
}

func deepseekCall(model *models.DeepseekModel, response *types.CurlResponse, curl *models.CurlRequest) ([]byte, error) {
	assistantContent, err := response.GetAssistantContent()
	if err != nil {
		return nil, err
	}

	properties := curl.GetOllamaSchemaProperties()
	properties["IsCorrect"] = ollama.OllamaFormatProperty{
		Type:        "boolean",
		Description: "This holds the true or false value for the Statement",
	}
	ollamaReq := ollama.OllamaRequest{
		Model: model.ModelName,
		Messages: []*ollama.OllamaMessage{
			{
				Role:    "assistant",
				Content: assistantContent,
			},
			{
				Role:    "user",
				Content: curl.Body,
			},
		},
		Stream: false,
		Format: &ollama.OllamaFormat{
			Type:       "object",
			Properties: properties,
			Required: func() []string {
				requiredKeys := make([]string, 0, len(properties))
				for key := range properties {
					requiredKeys = append(requiredKeys, key)
				}
				return requiredKeys
			}(),
		},
		Options: &ollama.OllamaOptions{
			Temperature: 0.5,
		},
		Think: true,
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/chat", model.BaseURL)
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
