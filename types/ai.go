package types

import (
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"fmt"

	"gorm.io/gorm"

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
	ID        uint   `json:"-"`
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
		rules = append(rules, validation.Field(&r.BaseURL, validation.Length(0, 500))) // Optional BaseURL
	}

	return validation.ValidateStruct(r, rules...)
}

func (dr *AIModelRequest) ToModel() (models.AIModelInterface, error) {
	err := dr.Validate()
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrInvalidRequestBody, err)
	}
	aiModel := models.AIModel{
		Model: gorm.Model{
			ID: dr.ID,
		},
		Type:    dr.Type,
		BaseURL: &dr.BaseURL,
	}
	switch dr.Type {
	case "deepseek", "local":
		return &models.DeepseekModel{
			AIModel:   aiModel,
			Name:      dr.Name,
			ModelName: dr.ModelName,
			Size:      dr.Size,
		}, nil
	case "gemini":
		geminiModel := &models.GeminiModel{
			AIModel:   aiModel,
			Name:      dr.Name,
			ModelName: dr.ModelName,
			APISecret: models.EncryptedString(dr.APISecret),
		}
		return geminiModel, nil
	case "openai":
		openaiModel := &models.OpenAIModel{
			AIModel:   aiModel,
			Name:      dr.Name,
			ModelName: dr.ModelName,
			APISecret: models.EncryptedString(dr.APISecret),
		}
		return openaiModel, nil
	default:
		return nil, errutil.NewAppError(errutil.ErrUnsupportedAIModelType, fmt.Errorf("unsupported AI model type: %s", dr.Type))
	}
}
