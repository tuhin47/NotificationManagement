package models

import (
	"gorm.io/gorm"
)

type UserLLM struct {
	gorm.Model
	RequestID uint     `gorm:"index"`
	IsActive  bool     `gorm:"default:true"`
	AiModelID uint     // Foreign key for AIModel
	AiModel   *AIModel `gorm:"foreignKey:AiModelID"`
	// ModelName string `gorm:"size:255;not null;index"`
	// Type     string `gorm:"size:50;check:type IN ('local','openai','gemini')"`
	//Parameters   JSON   `gorm:"type:jsonb"`
}
