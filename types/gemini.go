package types

import (
	"NotificationManagement/models"
)

type GeminiModelResponse struct {
	ID         uint   `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	ModelName  string `json:"model"`
	APISecret  string `json:"api_secret"`
	ModifiedAt string `json:"modified_at"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func FromGeminiModel(model *models.GeminiModel) *GeminiModelResponse {
	return &GeminiModelResponse{
		ID:        model.ID,
		Type:      model.Type,
		Name:      model.Name,
		ModelName: model.ModelName,
		APISecret: model.APISecret,
		CreatedAt: model.CreatedAt.Format(ResponseDateFormat),
		UpdatedAt: model.UpdatedAt.Format(ResponseDateFormat),
	}
}
