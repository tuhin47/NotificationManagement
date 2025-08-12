package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"context"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	domain.Repository[models.User, uint]
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepositoryImpl{
		Repository: NewSQLRepository[models.User](db),
	}
}

func (r *UserRepositoryImpl) FindByKeycloakID(keycloakID string, ctx context.Context) (*models.User, error) {
	var user models.User

	err := r.GetDB(ctx).Where("keycloak_id = ?", keycloakID).First(&user).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return &user, nil
}
