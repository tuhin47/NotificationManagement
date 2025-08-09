package models

import (
	"gorm.io/gorm"
)

type Telegram struct {
	gorm.Model
	ChatID int64  `gorm:"uniqueIndex;not null"`
	UserID *uint  `gorm:"index"`
	User   User   `gorm:"foreignKey:UserID"`
	Otp    string `gorm:"size:255;not null"`
}
