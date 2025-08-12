package domain

import (
	"NotificationManagement/models"
	"context"
)

type UserRepository interface {
	Repository[models.User, uint]
	FindByKeycloakID(keycloakID string, ctx context.Context) (*models.User, error)
}

type UserService interface {
	CommonService[models.User]
	RegisterOrUpdateUser(ctx context.Context, user *models.User) (*models.User, error)
}

type UserController interface {
}
