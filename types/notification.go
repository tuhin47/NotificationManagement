package types

import "NotificationManagement/models"

type Notification struct {
	Subject  string
	Message  string
	Channels []string
	UserId   uint
	User     *models.User
}

type NotifyRequest struct {
	UserId   uint     `json:"user_id"`
	Subject  string   `json:"subject"`
	Message  string   `json:"message"`
	Channels []string `json:"channels"`
}
