package models

import (
	"gorm.io/gorm"
)

type AIModelInterface interface {
	GetType() string
}

type AIModel struct {
	gorm.Model
	Type string `gorm:"size:10;check:type IN ('local','openai','gemini','deepseek')"`
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
	Name      string `gorm:"size:255;not null" json:"name"`
	ModelName string `gorm:"size:255;not null;check:model_name <> '';index:idx_ai_model_model_secret,unique" json:"model"`
	APISecret string `gorm:"size:500;index:idx_ai_model_model_secret,unique" json:"api_secret"`
}

func (d *AIModel) GetType() string {
	return d.Type
}

func (d *GeminiModel) GetType() string {
	return d.Type
}

func (d *DeepseekModel) GetType() string {
	return d.Type
}

func (*DeepseekModel) TableName() string {
	return "ai_models"
}

func (*GeminiModel) TableName() string {
	return "ai_models"
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
