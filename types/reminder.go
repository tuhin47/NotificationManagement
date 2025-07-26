package types

import (
	"NotificationManagement/models"
	"time"
)

type ReminderRequest struct {
	RequestID       uint      `json:"request_id" validate:"required"`
	Message         string    `json:"message" validate:"required"`
	TriggeredTime   time.Time `json:"triggered_time" validate:"required"`
	NextTriggerTime time.Time `json:"next_trigger_time" validate:"required"`
	Occurrence      uint      `json:"occurrence"`
	Recurrence      string    `json:"recurrence" validate:"required,oneof=once minutes hour daily weekly"`
}

type ReminderResponse struct {
	ID              uint      `json:"id"`
	RequestID       uint      `json:"request_id"`
	Message         string    `json:"message"`
	TriggeredTime   time.Time `json:"triggered_time"`
	NextTriggerTime time.Time `json:"next_trigger_time"`
	Occurrence      uint      `json:"occurrence"`
	Recurrence      string    `json:"recurrence"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

// ToModel converts a types.ReminderRequest to a models.Reminder
func (rr *ReminderRequest) ToModel() *models.Reminder {
	return &models.Reminder{
		RequestID:       rr.RequestID,
		Message:         rr.Message,
		TriggeredTime:   rr.TriggeredTime,
		NextTriggerTime: rr.NextTriggerTime,
		Occurrence:      rr.Occurrence,
		Recurrence:      rr.Recurrence,
	}
}

// FromModel converts a models.Reminder to a types.ReminderResponse
func FromReminderModel(model *models.Reminder) *ReminderResponse {
	return &ReminderResponse{
		ID:              model.ID,
		RequestID:       model.RequestID,
		Message:         model.Message,
		TriggeredTime:   model.TriggeredTime,
		NextTriggerTime: model.NextTriggerTime,
		Occurrence:      model.Occurrence,
		Recurrence:      model.Recurrence,
		CreatedAt:       model.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       model.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
