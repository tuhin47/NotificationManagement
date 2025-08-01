package types

import "NotificationManagement/models"

type MakeAIRequestPayload struct {
	CurlRequestID uint `json:"curl_request_id" validate:"required"`
	ModelID       uint `json:"model_id" validate:"required"`
}

type AIModelRequest struct {
	Name       string `json:"name" validate:"required"`
	Type       string `json:"type" validate:"required"`
	ModelName  string `json:"model" validate:"required"`
	BaseURL    string `json:"base_url"`
	ModifiedAt string `json:"modified_at" validate:"required"`
	Size       int64  `json:"size" validate:"required"`
}

// ToModel converts a types.AIModelRequest to a models.DeepseekModel
func (dr *AIModelRequest) ToModel() *models.DeepseekModel {
	m := &models.AIModel{
		Type: "local",
	}
	return &models.DeepseekModel{
		AIModel:    *m,
		Name:       dr.Name,
		ModelName:  dr.ModelName,
		BaseURL:    dr.BaseURL,
		ModifiedAt: dr.ModifiedAt,
		Size:       dr.Size,
	}
}

func FromDeepseekModel(model *models.DeepseekModel) *DeepseekModelResponse {
	return &DeepseekModelResponse{
		ID:         model.ID,
		Name:       model.Name,
		ModelName:  model.ModelName,
		BaseURL:    model.BaseURL,
		ModifiedAt: model.ModifiedAt,
		Size:       model.Size,
		CreatedAt:  model.CreatedAt.Format(ResponseDateFormat),
		UpdatedAt:  model.UpdatedAt.Format(ResponseDateFormat),
	}
}
