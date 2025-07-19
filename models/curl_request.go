package models

import (
	"gorm.io/gorm"
)

type CurlRequest struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	URL       string         `gorm:"type:text" json:"url"`
	Method    string         `gorm:"type:varchar(10)" json:"method"`
	Headers   string         `gorm:"type:text" json:"headers"` // JSON string
	Body      string         `gorm:"type:text" json:"body"`
	RawCurl   string         `gorm:"type:text" json:"rawCurl"`
	CreatedAt int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
