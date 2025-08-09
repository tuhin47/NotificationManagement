package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system, linked to Keycloak.
type User struct {
	gorm.Model
	KeycloakID string      `gorm:"uniqueIndex;not null"`
	Username   string      `gorm:"size:255;not null"`
	Email      string      `gorm:"size:255;uniqueIndex;not null"`
	Roles      string      `gorm:"type:text"`
	Telegram   *[]Telegram `gorm:"foreignKey:UserID"`
}

func (u *User) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*User); ok {
		copyFields(u, src)
	}
}
