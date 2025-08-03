package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/repositories"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"context"
	"google.golang.org/genai"
)

type GeminiServiceImpl struct {
	domain.CommonService[models.GeminiModel]
	Repo        domain.GeminiModelRepository
	CurlService domain.CurlService
}

func NewGeminiService(repo domain.GeminiModelRepository, curlService domain.CurlService) domain.GeminiService {
	service := &GeminiServiceImpl{
		Repo:        repo,
		CurlService: curlService,
	}
	service.CommonService = NewCommonService[models.GeminiModel](repo, service)
	return service
}

func (s *GeminiServiceImpl) GetContext() context.Context {
	background := context.Background()
	f := []repositories.Filter{
		{"type", "=", "gemini"},
	}
	return context.WithValue(background, repositories.ContextStruct{}, &repositories.ContextStruct{Filter: &f})
}

func (s *GeminiServiceImpl) MakeAIRequest(mod *models.AIModel, requestId uint) (*genai.GenerateContentResponse, error) {

	curl, err := s.CurlService.GetModelByID(requestId)
	if err != nil {
		return nil, err
	}
	curlResponse, err := s.CurlService.ExecuteCurl(curl)
	if err != nil {
		return nil, err
	}
	model, err := s.GetModelByID(requestId)
	if err != nil {
		return nil, err
	}
	respBody, err := geminiCall(model, curlResponse, curl)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	return respBody, nil
}

func geminiCall(model *models.GeminiModel, response *types.CurlResponse, req *models.CurlRequest) (*genai.GenerateContentResponse, error) {
	assistantContent, err := response.GetAssistantContent()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	// The client gets the API key from the environment variable `GEMINI_API_KEY`.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: model.APISecret,
	})
	if err != nil {
		return nil, err
	}
	gr := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: assistantContent},
			},
		},
		{
			Role: "model",
			Parts: []*genai.Part{
				{Text: "Please check the current rate from the json.Is it greater than 125 ? Return Json Response "},
			},
		},
	}
	properties := map[string]*genai.Schema{
		"IsCorrect": {Type: genai.TypeBoolean},
		"CurrentRate": {
			Type:        genai.TypeNumber,
			Description: "description",
		},
	}
	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
		},
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:       genai.TypeObject,
			Properties: properties,
		},
	}
	result, err := client.Models.GenerateContent(
		ctx,
		model.ModelName,
		gr,
		config,
	)
	if err != nil {
		return nil, err
	}
	logger.Info(result.Text())

	return result, err
}
