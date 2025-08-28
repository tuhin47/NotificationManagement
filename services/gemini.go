package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/repositories"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"context"
	"encoding/json"
	"os"

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

func (s *GeminiServiceImpl) ProcessContext(ctx context.Context) context.Context {
	if txContext, ok := repositories.GetTxContext(ctx); ok {
		filters := append(txContext.Filter, repositories.NewFilter("type", "=", s.GetModelType()))
		txContext.Filter = filters
	}
	return ctx
}

func (s *GeminiServiceImpl) MakeAIRequest(c context.Context, m *models.AIModel, requestId uint) (interface{}, error) {

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
	respBody, err := geminiCall(c, model, curlResponse, curl)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}
	return respBody, nil
}

func (s *GeminiServiceImpl) GetAIJsonResponse(c context.Context, m *models.AIModel, requestId uint) (map[string]interface{}, error) {
	request, err := s.MakeAIRequest(c, m, requestId)
	if err != nil {
		return nil, err
	}
	resp, _ := request.(*genai.GenerateContentResponse)
	var aiResp map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Text()), &aiResp); err != nil {
		return nil, err
	}
	return aiResp, nil
}

func (s *GeminiServiceImpl) GetModelType() string {
	return "gemini"
}

func geminiCall(ctx context.Context, model *models.GeminiModel, response *types.CurlResponse, req *models.CurlRequest) (*genai.GenerateContentResponse, error) {
	assistantContent, err := response.GetAssistantContent(req.ResponseType)
	if err != nil {
		return nil, err
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: model.GetAPIKey(),
		HTTPOptions: genai.HTTPOptions{
			BaseURL: model.GetBaseURL(),
		},
	})
	if err != nil {
		return nil, err
	}
	var parts []*genai.Part

	if req.ResponseType == types.ResponseTypeHTML {
		fileContent, err := os.ReadFile(*assistantContent)
		if err != nil {
			return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
		}
		parts = append(parts, &genai.Part{
			InlineData: &genai.Blob{
				MIMEType: "text/html", // Assuming it's HTML for /tmp/1.html
				Data:     fileContent,
			},
		})
	} else {
		parts = append(parts, &genai.Part{Text: *assistantContent})
	}

	parts = append(parts, &genai.Part{Text: req.Body})

	gr := []*genai.Content{
		{
			Role:  genai.RoleModel,
			Parts: parts,
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

func (s *GeminiServiceImpl) CreateAIModel(c context.Context, model any) error {
	geminiModel := (model).(*models.GeminiModel)
	return s.CommonService.CreateModel(c, geminiModel)
}

func (s *GeminiServiceImpl) UpdateAIModel(c context.Context, model any) (any, error) {
	geminiModel := (model).(*models.GeminiModel)
	return s.CommonService.UpdateModel(c, geminiModel.ID, geminiModel)
}

func (s *GeminiServiceImpl) GetAIModelById(ctx context.Context, id uint) (any, error) {
	return s.GetModelById(ctx, id, nil)
}

func (s *GeminiServiceImpl) GetAllAIModels(ctx context.Context) ([]any, error) {
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
