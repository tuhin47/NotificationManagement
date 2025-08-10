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

func (s *DeepseekServiceImpl) MakeAIRequest(c context.Context, m *models.AIModel, requestId uint) (interface{}, error) {
	curl, err := s.CurlService.GetModelById(c, requestId, nil)
	if err != nil {
		return nil, err
	}
	curlResponse, err := s.CurlService.ProcessCurlRequest(c, curl)
	if err != nil {
		return nil, err
	}
	model, err := s.GetModelById(c, m.ID, nil)
	if err != nil {
		return nil, err
	}
	respBody, err := deepseekCall(model, curlResponse, curl)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	var ollamaResp ollama.Response
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	a := any(ollamaResp)
	return &a, nil
}

func (s *DeepseekServiceImpl) GetAIJsonResponse(c context.Context, m *models.AIModel, requestId uint) (map[string]interface{}, error) {
	request, err := s.MakeAIRequest(c, m, requestId)
	if err != nil {
		return nil, err
	}
	resp, _ := request.(*ollama.Response)

	aiResp := make(map[string]interface{})
	if err := json.Unmarshal([]byte(resp.Message.Content), &aiResp); err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}
	return aiResp, nil

}

func (s *DeepseekServiceImpl) PullModel(_ context.Context, model *models.DeepseekModel) error {
	// Implementation for pulling/downloading the model
	payload := map[string]interface{}{
		"name": model.ModelName,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return errutil.NewAppError(errutil.ErrAIMarshalRequestFailed, err)
	}

	url := fmt.Sprintf("%s/api/pull", model.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errutil.NewAppError(errutil.ErrAICreateRequestFailed, err)
	}

	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return errutil.NewAppError(errutil.ErrAIPullModelFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errutil.NewAppError(errutil.ErrAIPullModelFailed, fmt.Errorf("status code: %d", resp.StatusCode))
	}

	return nil
}

func deepseekCall(model *models.DeepseekModel, response *types.CurlResponse, curl *models.CurlRequest) ([]byte, error) {
	assistantContent, err := response.GetAssistantContent(curl.ResponseType)
	if err != nil {
		return nil, err
	}

	properties := curl.GetOllamaSchemaProperties()
	properties["IsCorrect"] = ollama.FormatProperty{
		Type:        "boolean",
		Description: "This holds the true or false value for the Statement",
	}
	if curl.ResponseType == types.ResponseTypeHTML {
		s := response.Body.(string)
		assistantContent = &s
	}

	ollamaReq := ollama.Request{
		Model: model.ModelName,
		Messages: []*ollama.Message{
			{
				Role:    "assistant",
				Content: *assistantContent,
			},
			{
				Role:    "user",
				Content: curl.Body,
			},
		},
		Stream: false,
		Format: &ollama.Format{
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
		Options: &ollama.Options{
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

func (s *DeepseekServiceImpl) GetModelType() string {
	return "deepseek"
}

func (s *DeepseekServiceImpl) CreateAIModel(c context.Context, model any) error {
	deepseekModel := (model).(*models.DeepseekModel)
	return s.CreateModel(c, deepseekModel)
}
func (s *DeepseekServiceImpl) UpdateAIModel(c context.Context, model any) (any, error) {
	deepseekModel := (model).(*models.DeepseekModel)
	return s.UpdateModel(c, deepseekModel.ID, deepseekModel)
}

func (s *DeepseekServiceImpl) GetAIModelById(ctx context.Context, id uint) (any, error) {
	return s.GetModelById(ctx, id, nil)
}
func (s *DeepseekServiceImpl) GetAllAIModels(ctx context.Context) ([]any, error) {
	allModels, err := s.GetAllModels(ctx, 100, 0)
	if err != nil {
		return nil, err
	}
	i := make([]any, len(allModels))
	for idx, model := range allModels {
		i[idx] = model
	}
	return i, err
}
