package models

import (
	"NotificationManagement/config"

	"gorm.io/gorm"
)

type AIModelInterface interface {
	GetType() string
}

type AIModel struct {
	gorm.Model
	Type string `gorm:"size:10;check:type IN ('local','openai','gemini','deepseek')"`
}

type OpenAIModel struct {
	AIModel   `mapper:"inherit"`
	Name      string          `gorm:"size:255;not null" json:"name"`
	ModelName string          `gorm:"size:255;not null;check:model_name <> '';index:idx_ai_model_model_secret,unique" json:"model"`
	APISecret EncryptedString `gorm:"size:500;index:idx_ai_model_model_secret,unique" json:"-"`
}

type DeepseekModel struct {
	AIModel   `mapper:"inherit"`
	Name      string `gorm:"size:255;not null" json:"name"`
	ModelName string `gorm:"size:255;not null;check:model_name <> '';index:idx_ai_model_model_url,unique" json:"model"`
	BaseURL   string `gorm:"size:500;index:idx_ai_model_model_url,unique" json:"base_url"`
	Size      int64  `json:"size"`
}

type GeminiModel struct {
	AIModel   `mapper:"inherit"`
	Name      string          `gorm:"size:255;not null" json:"name"`
	ModelName string          `gorm:"size:255;not null;check:model_name <> '';index:idx_ai_model_model_secret,unique" json:"model"`
	APISecret EncryptedString `gorm:"size:500;index:idx_ai_model_model_secret,unique" json:"-"`
}

func (d *AIModel) GetType() string {
	return d.Type
}

func (*AIModel) TableName() string {
	return "ai_models"
}

func (d *GeminiModel) GetAPIKey() string {
	if config.IsDevelopment() && config.Development().GeminiKey != "" {
		return config.Development().GeminiKey
	}
	return string(d.APISecret)
}

func (d *OpenAIModel) GetAPIKey() string {
	if config.IsDevelopment() && config.Development().OpenAIKey != "" {
		return config.Development().OpenAIKey
	}
	return string(d.APISecret)
}

func (d *DeepseekModel) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*DeepseekModel); ok {
		copyFields(d, src)
	}
}
func (d *GeminiModel) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*GeminiModel); ok {
		copyFields(d, src)
	}
}

func (d *OpenAIModel) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*OpenAIModel); ok {
		copyFields(d, src)
	}
}
