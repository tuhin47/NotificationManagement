package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"context"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	domain.Repository[models.User, uint]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepositoryImpl{
		Repository: NewSQLRepository[models.User](db),
		db:         db,
	}
}

func (r *UserRepositoryImpl) FindByKeycloakID(keycloakID string, ctx context.Context) (*models.User, error) {
	var user models.User
	err := r.db.Where("keycloak_id = ?", keycloakID).First(&user).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return &user, nil
}
