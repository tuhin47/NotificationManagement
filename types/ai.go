package types

import (
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"fmt"
)

type MakeAIRequestPayload struct {
	CurlRequestID uint `json:"curl_request_id" validate:"required"`
	ModelID       uint `json:"model_id" validate:"required"`
}

type AIModelRequest struct {
	Name      string `json:"name" validate:"required"`
	Type      string `json:"type" validate:"required"`
	ModelName string `json:"model" validate:"required"`
	BaseURL   string `json:"base_url"`
	APISecret string `json:"api_secret"`
	Size      int64  `json:"size"`
}

// ToModel converts a types.AIModelRequest to a models.DeepseekModel or models.GeminiModel
func (dr *AIModelRequest) ToModel() (interface{}, error) {
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
