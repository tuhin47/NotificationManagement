package types

import (
	"NotificationManagement/models"
)

type DeepseekModelRequest struct {
	Name       string `json:"name" validate:"required"`
	ModelName  string `json:"model" validate:"required"`
	ModifiedAt string `json:"modified_at" validate:"required"`
	Size       int64  `json:"size" validate:"required"`
}

type DeepseekModelResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	ModelName  string `json:"model"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ToModel converts a types.DeepseekModelRequest to a models.DeepseekModel
func (dr *DeepseekModelRequest) ToModel() *models.DeepseekModel {
	m := &models.AIModel{
		Type: "local",
	}
	return &models.DeepseekModel{
		AIModel:    *m,
		Name:       dr.Name,
		ModelName:  dr.ModelName,
		ModifiedAt: dr.ModifiedAt,
		Size:       dr.Size,
	}
}

func FromDeepseekModel(model *models.DeepseekModel) *DeepseekModelResponse {
	return &DeepseekModelResponse{
		ID:         model.ID,
		Name:       model.Name,
		ModelName:  model.ModelName,
		ModifiedAt: model.ModifiedAt,
		Size:       model.Size,
		CreatedAt:  model.CreatedAt.Format(ResponseDateFormat),
		UpdatedAt:  model.UpdatedAt.Format(ResponseDateFormat),
	}
}
