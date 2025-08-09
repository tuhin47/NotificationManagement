package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"context"
	"errors"
)

type UserServiceImpl struct {
	domain.CommonService[models.User]
	UserRepo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	service := &UserServiceImpl{
		UserRepo: repo,
	}
	service.CommonService = NewCommonService(repo, service)
	return service
}

func (s *UserServiceImpl) RegisterOrUpdateUser(user *models.User) (*models.User, error) {
	ctx := context.Background()
	existingUser, err := s.UserRepo.FindByKeycloakID(user.KeycloakID, ctx)
	if err != nil {
		var appErr *errutil.AppError
		if errors.As(err, &appErr) && appErr.Code == errutil.ErrRecordNotFound {
			// User not found, create new user
			err = s.UserRepo.Create(ctx, user)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}

	// User found, update existing user
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.Roles = user.Roles
	err = s.UserRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}
	return existingUser, nil
}

func (s *UserServiceImpl) GetContext() context.Context {
	return context.Background()
}
