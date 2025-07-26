package models

import (
	"gorm.io/gorm"
)

type UserLLM struct {
	gorm.Model
	RequestID uint   `gorm:"index"`
	ModelName string `gorm:"size:255;not null;index"`
	Type      string `gorm:"size:50;check:type IN ('local','openai','gemini')"`
	IsActive  bool   `gorm:"default:true"`
	//Parameters   JSON   `gorm:"type:jsonb"`
}
