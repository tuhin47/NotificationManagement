package domain

import (
	"NotificationManagement/models"
	"context"
)

// UserRepository is an interface for user data operations.
type UserRepository interface {
	Repository[models.User, uint]
	FindByKeycloakID(keycloakID string, ctx context.Context) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}

// UserService is an interface for user-related business logic.
type UserService interface {
	CommonService[models.User]
	RegisterOrUpdateUser(user *models.User) (*models.User, error)
}

// UserController is an interface for user HTTP handlers.
type UserController interface {
}
