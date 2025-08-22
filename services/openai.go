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

	"github.com/sashabaranov/go-openai"
)

// JSONSchemaProperty represents a property in a JSON schema
type JSONSchemaProperty struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// JSONSchema represents a JSON schema structure
type JSONSchema struct {
	Type       string                        `json:"type"`
	Properties map[string]JSONSchemaProperty `json:"properties"`
	Required   []string                      `json:"required"`
}

// MarshalJSON implements json.Marshaler interface
func (j JSONSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string                        `json:"type"`
		Properties map[string]JSONSchemaProperty `json:"properties"`
		Required   []string                      `json:"required"`
	}{
		Type:       j.Type,
		Properties: j.Properties,
		Required:   j.Required,
	})
}

type OpenAIServiceImpl struct {
	domain.CommonService[models.OpenAIModel]
	CurlService domain.CurlService
}

func NewOpenAIService(repo domain.OpenAIModelRepository, curl domain.CurlService) domain.OpenAIService {
	service := &OpenAIServiceImpl{
		CurlService: curl,
	}
	service.CommonService = NewCommonService(repo, service)
	return service
}

func (s *OpenAIServiceImpl) ProcessContext(ctx context.Context) context.Context {
	if txContext, ok := repositories.GetTxContext(ctx); ok {
		filters := append(txContext.Filter, repositories.NewFilter("type", "=", s.GetModelType()))
		txContext.Filter = filters
	}
	return ctx
}

func (s *OpenAIServiceImpl) MakeAIRequest(c context.Context, m *models.AIModel, requestId uint) (interface{}, error) {
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
	respBody, err := openAICall(c, model, curlResponse, curl)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func (s *OpenAIServiceImpl) GetAIJsonResponse(c context.Context, m *models.AIModel, requestId uint) (map[string]interface{}, error) {
	request, err := s.MakeAIRequest(c, m, requestId)
	if err != nil {
		return nil, err
	}
	resp, _ := request.(*openai.ChatCompletionResponse)
	var aiResp map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &aiResp); err != nil {
		return map[string]interface{}{
			"comment": resp.Choices[0].Message.Content,
		}, nil
	}
	return aiResp, nil
}

func (s *OpenAIServiceImpl) GetModelType() string {
	return "openai"
}

// createJSONSchema creates a JSON schema from the CurlRequest additional fields
func createJSONSchema(req *models.CurlRequest) JSONSchema {
	properties := make(map[string]JSONSchemaProperty)
	required := []string{"IsCorrect"}

	properties["IsCorrect"] = JSONSchemaProperty{
		Type:        "boolean",
		Description: "Indicates whether the response is correct",
	}

	// Add properties from additional fields
	if req.AdditionalFields != nil {
		for _, field := range *req.AdditionalFields {
			var schemaType string
			switch field.Type {
			case "number":
				schemaType = "number"
			case "boolean":
				schemaType = "boolean"
			default:
				schemaType = "string"
			}

			properties[field.PropertyName] = JSONSchemaProperty{
				Type:        schemaType,
				Description: field.Description,
			}
			required = append(required, field.PropertyName)
		}
	}

	return JSONSchema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

func openAICall(ctx context.Context, model *models.OpenAIModel, response *types.CurlResponse, req *models.CurlRequest) (*openai.ChatCompletionResponse, error) {
	assistantContent, err := response.GetAssistantContent(req.ResponseType)
	if err != nil {
		return nil, err
	}

	config := openai.DefaultConfig(model.GetAPIKey())
	config.BaseURL = "https://api.closerouter.com/v1"
	client := openai.NewClientWithConfig(config)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleAssistant,
			Content: *assistantContent,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: req.Body,
		},
	}

	// Create JSON schema from additional fields for structured output
	jsonSchema := createJSONSchema(req)

	logger.Debug(jsonSchema.Type)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model.ModelName,
			Messages: messages,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:        "response_schema",
					Description: "Structured response with required fields",
					Schema:      jsonSchema,
					Strict:      true,
				},
			},
		},
	)

	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	return &resp, nil
}

func (s *OpenAIServiceImpl) CreateAIModel(c context.Context, model any) error {
	openaiModel := (model).(*models.OpenAIModel)
	return s.CreateModel(c, openaiModel)
}

func (s *OpenAIServiceImpl) UpdateAIModel(c context.Context, model any) (any, error) {
	openaiModel := (model).(*models.OpenAIModel)
	return s.UpdateModel(c, openaiModel.ID, openaiModel)
}

func (s *OpenAIServiceImpl) GetAIModelById(ctx context.Context, id uint) (any, error) {
	return s.GetModelById(ctx, id, nil)
}

func (s *OpenAIServiceImpl) GetAllAIModels(ctx context.Context) ([]any, error) {
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
