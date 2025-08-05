package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system, linked to Keycloak.
type User struct {
	gorm.Model
	KeycloakID string `gorm:"uniqueIndex;not null"` // Keycloak's internal user ID
	Username   string `gorm:"size:255;not null"`
	Email      string `gorm:"size:255;uniqueIndex;not null"`
	Roles      string `gorm:"type:text"` // Comma-separated roles from Keycloak
}

func (u *User) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*User); ok {
		copyFields(u, src)
	}
}
