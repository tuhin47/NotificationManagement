package types

import (
	"NotificationManagement/models"
)

type LLMRequest struct {
	RequestID uint   `json:"request_id"`
	ModelName string `json:"model_name" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=local openai gemini"`
	IsActive  bool   `json:"is_active"`
}

type LLMResponse struct {
	ID        uint   `json:"id"`
	RequestID uint   `json:"request_id"`
	ModelName string `json:"model_name"`
	Type      string `json:"type"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ToModel converts a types.LLMRequest to a models.UserLLM
func (lr *LLMRequest) ToModel() *models.UserLLM {
	return &models.UserLLM{
		RequestID: lr.RequestID,
		ModelName: lr.ModelName,
		Type:      lr.Type,
		IsActive:  lr.IsActive,
	}
}

// FromModel converts a models.UserLLM to a types.LLMResponse
func FromLLMModel(model *models.UserLLM) *LLMResponse {
	return &LLMResponse{
		ID:        model.ID,
		RequestID: model.RequestID,
		ModelName: model.ModelName,
		Type:      model.Type,
		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: model.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
