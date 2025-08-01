package models

import (
	"gorm.io/gorm"
)

type AIModelInterface interface {
	UpdateFromModel(source AIModelInterface)
	GetBaseURL() string
}

type AIModel struct {
	gorm.Model
	Type string `gorm:"size:50;check:type IN ('local','openai','gemini','deepseek')"`
}

type DeepseekModel struct {
	AIModel
	Name       string `gorm:"size:255;not null" json:"name"`
	ModelName  string `gorm:"size:255;not null;check:model_name <> ''" json:"model"`
	BaseURL    string `gorm:"size:500" json:"base_url"`
	ModifiedAt string `gorm:"size:50" json:"modified_at"`
	Size       int64  `gorm:"not null" json:"size"`
}

func (d *DeepseekModel) GetBaseURL() string {
	return d.BaseURL
}

func (*DeepseekModel) TableName() string {
	return "ai_models"
}

type GeminiModel struct {
	AIModel
	Name       string `gorm:"size:255;not null" json:"name"`
	ModelName  string `gorm:"size:255;not null;check:model_name <> ''" json:"model"`
	APISecret  string `gorm:"size:500" json:"api_secret"`
	ModifiedAt string `gorm:"size:50" json:"modified_at"`
	Size       int64  `gorm:"not null" json:"size"`
}

func (*GeminiModel) GetBaseURL() string {
	return "https://api.gemini.com"
}

func (*GeminiModel) TableName() string {
	return "ai_models"
}

func (d *DeepseekModel) UpdateFromModel(source AIModelInterface) {
	if src, ok := source.(*DeepseekModel); ok {
		d.Name = src.Name
		d.ModelName = src.ModelName
		d.BaseURL = src.BaseURL
		d.ModifiedAt = src.ModifiedAt
		d.Size = src.Size
	}
}
func (g *GeminiModel) UpdateFromModel(source AIModelInterface) {
	if src, ok := source.(*GeminiModel); ok {
		g.Name = src.Name
		g.ModelName = src.ModelName
		g.APISecret = src.APISecret
		g.ModifiedAt = src.ModifiedAt
		g.Size = src.Size
	}
}
