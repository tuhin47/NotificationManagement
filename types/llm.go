package types

import (
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type LLMRequest struct {
	RequestID uint `json:"request_id"`
	AIModelID uint `json:"ai_model_id"`
	IsActive  bool `json:"is_active"`
}

func (r *LLMRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.RequestID, validation.Required),
		validation.Field(&r.AIModelID, validation.Required),
	)
}

type LLMResponse struct {
	ID        uint   `json:"id"`
	RequestID uint   `json:"request_id"`
	AIModelID uint   `json:"ai_model_id"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ToModel converts a types.LLMRequest to a models.UserLLM
func (lr *LLMRequest) ToModel() (*models.UserLLM, error) {
	err := lr.Validate()
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrInvalidRequestBody, err)
	}
	return &models.UserLLM{
		RequestID: lr.RequestID,
		AiModelID: lr.AIModelID,
		IsActive:  lr.IsActive,
	}, nil
}

// FromModel converts a models.UserLLM to a types.LLMResponse
func FromLLMModel(model *models.UserLLM) *LLMResponse {
	return &LLMResponse{
		ID:        model.ID,
		RequestID: model.RequestID,
		AIModelID: model.AiModelID,
		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt.Format(ResponseDateFormat),
		UpdatedAt: model.UpdatedAt.Format(ResponseDateFormat),
	}
}
