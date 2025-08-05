package controllers

import (
	"NotificationManagement/domain"
)

type UserControllerImpl struct {
	UserService domain.UserService
}

func NewUserController(userService domain.UserService) domain.UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}
