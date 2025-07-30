package models

import (
	"gorm.io/gorm"
)

type CurlRequest struct {
	gorm.Model
	URL                    string                  `gorm:"type:text" json:"url"`
	Method                 string                  `gorm:"type:varchar(10)" json:"method"`
	Headers                string                  `gorm:"type:text" json:"headers"`
	Body                   string                  `gorm:"type:text" json:"body"`
	RawCurl                string                  `gorm:"type:text" json:"rawCurl"`
	Reminders              *[]Reminder             `gorm:"foreignKey:RequestID"`
	LLMs                   *[]UserLLM              `gorm:"foreignKey:RequestID"`
	OllamaFormatProperties *[]OllamaFormatProperty `gorm:"foreignKey:RequestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ollama_format_properties"`
}

type OllamaFormatProperty struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	PropertyName string `gorm:"type:varchar(100)" json:"property_name"`
	Type         string `gorm:"type:varchar(10)" json:"type"` // allowed: number, boolean, text
	Description  string `gorm:"type:text" json:"description,omitempty"`
	RequestID    uint   `json:"request_id"`
}
