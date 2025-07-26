package models

import (
	"gorm.io/gorm"
)

type UserLLM struct {
	gorm.Model
	RequestID uint     `gorm:"index:idx_request_ai_model,unique"`
	IsActive  bool     `gorm:"default:true"`
	AiModelID uint     `gorm:"index:idx_request_ai_model,unique"` // Foreign key for AIModel
	AiModel   *AIModel `gorm:"foreignKey:AiModelID"`
	// ModelName string `gorm:"size:255;not null;index"`
	// Type     string `gorm:"size:50;check:type IN ('local','openai','gemini')"`
	//Parameters   JSON   `gorm:"type:jsonb"`
}
