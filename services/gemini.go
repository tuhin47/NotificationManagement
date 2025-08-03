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
	CurlService domain.CurlService
}

func NewGeminiService(repo domain.GeminiModelRepository, curlService domain.CurlService) domain.GeminiService {
	service := &GeminiServiceImpl{
		CurlService: curlService,
	}
	service.CommonService = NewCommonService(repo, service)
	return service
}

func (s *GeminiServiceImpl) GetContext() context.Context {
	background := context.Background()
	f := []repositories.Filter{
		{Field: "type", Op: "=", Value: "gemini"},
	}
	return context.WithValue(background, repositories.ContextStruct{}, &repositories.ContextStruct{Filter: &f})
}

func (s *GeminiServiceImpl) MakeAIRequest(aiModel *models.AIModel, requestId uint) (interface{}, error) {

	curl, err := s.CurlService.GetModelByID(requestId)
	if err != nil {
		return nil, err
	}
	curlResponse, err := s.CurlService.ExecuteCurl(curl)
	if err != nil {
		return nil, err
	}
	model, err := s.GetModelByID(aiModel.ID)
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
			Role: genai.RoleModel,
			Parts: []*genai.Part{
				{Text: assistantContent},
			},
		},
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: req.Body},
			},
		},
	}
	properties := req.GetGenaiSchemaProperties()
	properties["IsCorrect"] = &genai.Schema{
		Type:        genai.TypeBoolean,
		Description: "The answer of the question",
	}

	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: false,
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
