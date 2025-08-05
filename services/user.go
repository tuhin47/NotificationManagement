package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"strings"
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

func (s *UserServiceImpl) GetInstance() domain.CommonService[models.User] {
	return s.CommonService
}

func (s *UserServiceImpl) CreateModel(c echo.Context, entity *models.User) error {
	return s.UserRepo.Create(s.GetContext(), entity)
}

func (s *UserServiceImpl) GetModelById(c echo.Context, id uint) (*models.User, error) {
	model, err := s.UserRepo.GetByID(s.GetContext(), id, nil)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *UserServiceImpl) GetAllModels(c echo.Context, limit, offset int) ([]models.User, error) {
	m, err := s.UserRepo.GetAll(s.GetContext(), limit, offset)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *UserServiceImpl) UpdateModel(c echo.Context, id uint, model *models.User) (*models.User, error) {
	existing, err := s.UserRepo.GetByID(s.GetContext(), id, nil)
	if err != nil {
		return nil, err
	}
	existing.Username = model.Username
	existing.Email = model.Email
	existing.Roles = model.Roles
	err = s.UserRepo.Update(s.GetContext(), existing)
	return existing, err
}

func (s *UserServiceImpl) DeleteModel(c echo.Context, id uint) error {
	return s.UserRepo.Delete(s.GetContext(), id)
}

func (s *UserServiceImpl) GetUserRoles(keycloakID string) ([]string, error) {
	ctx := context.Background()
	user, err := s.UserRepo.FindByKeycloakID(keycloakID, ctx)
	if err != nil {
		return nil, err
	}
	return strings.Split(user.Roles, ","), nil
}
