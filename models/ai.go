package models

import "gorm.io/gorm"

type AIModel struct {
	gorm.Model
	Type string `gorm:"size:50;check:type IN ('local','openai','gemini')"`
}

type DeepseekModel struct {
	AIModel
	Name       string `gorm:"size:255;not null" json:"name"`
	ModelName  string `gorm:"size:255;not null;uniqueIndex;check:model_name <> ''" json:"model"`
	ModifiedAt string `gorm:"size:50" json:"modified_at"`
	Size       int64  `gorm:"not null" json:"size"`
}

func (DeepseekModel) TableName() string {
	return "ai_models"
}
