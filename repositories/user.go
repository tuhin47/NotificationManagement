package repositories

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"

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

func (r *UserRepositoryImpl) FindByKeycloakID(keycloakID string) (*models.User, error) {
	var user models.User
	err := r.GetDB(nil).Where("keycloak_id = ?", keycloakID).First(&user).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.GetDB(nil).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, handleDbError(err)
	}
	return &user, nil
}
