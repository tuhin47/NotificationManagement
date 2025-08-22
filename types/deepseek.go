package types

import "NotificationManagement/models"

type DeepseekModelResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	ModelName string `json:"model"`
	BaseURL   string `json:"base_url"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func FromDeepseekModel(model *models.DeepseekModel) *DeepseekModelResponse {
	return &DeepseekModelResponse{
		ID:        model.ID,
		Name:      model.Name,
		Type:      model.Type,
		ModelName: model.ModelName,
		BaseURL:   model.GetBaseURL(),
		Size:      model.Size,
		CreatedAt: model.CreatedAt.Format(ResponseDateFormat),
		UpdatedAt: model.UpdatedAt.Format(ResponseDateFormat),
	}
}
