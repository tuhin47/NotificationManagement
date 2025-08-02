package types

import (
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type MakeAIRequestPayload struct {
	CurlRequestID uint `json:"curl_request_id"`
	ModelID       uint `json:"model_id"`
}

func (p *MakeAIRequestPayload) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.CurlRequestID, validation.Required),
		validation.Field(&p.ModelID, validation.Required),
	)
}

type AIModelRequest struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	ModelName string `json:"model"`
	BaseURL   string `json:"base_url"`
	APISecret string `json:"api_secret"`
	Size      int64  `json:"size"`
}

func (r *AIModelRequest) Validate() error {
	rules := []*validation.FieldRules{
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Type, validation.Required, validation.In("deepseek", "local", "openai", "gemini")),
		validation.Field(&r.ModelName, validation.Required, validation.Length(1, 255)),
	}

	// Add conditional validation based on Type
	switch r.Type {
	case "deepseek", "local":
		rules = append(rules, validation.Field(&r.BaseURL, validation.Required, validation.Length(1, 500)))
	case "gemini", "openai":
		rules = append(rules, validation.Field(&r.APISecret, validation.Required, validation.Length(1, 500)))

	}

	return validation.ValidateStruct(r, rules...)
}

// ToModel converts a types.AIModelRequest to a models.DeepseekModel or models.GeminiModel
func (dr *AIModelRequest) ToModel() (interface{}, error) {
	err := dr.Validate()
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrInvalidRequestBody, err)
	}
	aiModel := models.AIModel{
		Type: dr.Type,
	}
	switch dr.Type {
	case "deepseek", "local":
		return &models.DeepseekModel{
			AIModel:   aiModel,
			Name:      dr.Name,
			ModelName: dr.ModelName,
			BaseURL:   dr.BaseURL,
			Size:      dr.Size,
		}, nil
	case "gemini":
		return &models.GeminiModel{
			AIModel:   aiModel,
			Name:      dr.Name,
			ModelName: dr.ModelName,
			APISecret: dr.APISecret,
		}, nil
	default:
		return nil, errutil.NewAppError(errutil.ErrUnsupportedAIModelType, fmt.Errorf("unsupported AI model type: %s", dr.Type))
	}
}

func FromDeepseekModel(model *models.DeepseekModel) *DeepseekModelResponse {
	return &DeepseekModelResponse{
		ID:        model.ID,
		Name:      model.Name,
		Type:      model.Type,
		ModelName: model.ModelName,
		BaseURL:   model.BaseURL,
		Size:      model.Size,
		CreatedAt: model.CreatedAt.Format(ResponseDateFormat),
		UpdatedAt: model.UpdatedAt.Format(ResponseDateFormat),
	}
}
