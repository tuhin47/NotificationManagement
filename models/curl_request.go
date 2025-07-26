package models

import (
	"gorm.io/gorm"
)

type CurlRequest struct {
	gorm.Model
	URL       string      `gorm:"type:text" json:"url"`
	Method    string      `gorm:"type:varchar(10)" json:"method"`
	Headers   string      `gorm:"type:text" json:"headers"`
	Body      string      `gorm:"type:text" json:"body"`
	RawCurl   string      `gorm:"type:text" json:"rawCurl"`
	Reminders []*Reminder `gorm:"foreignKey:RequestID"`
	LLMs      []*UserLLM  `gorm:"foreignKey:RequestID"`
}
