package types

import (
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ReminderRequest struct {
	RequestID     uint       `json:"request_id"`
	AfterEvery    uint       `json:"after_every"`
	Message       string     `json:"message"`
	TriggeredTime time.Time  `json:"triggered_time"`
	Occurrence    uint       `json:"occurrence"`
	Recurrence    string     `json:"recurrence"`
	Upto          *time.Time `json:"upto,omitempty"`
}

func (r *ReminderRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.RequestID, validation.Required),
		validation.Field(&r.Message, validation.Required, validation.Length(1, 2048)),
		validation.Field(&r.TriggeredTime, validation.Required),
		validation.Field(&r.AfterEvery, validation.When(r.Recurrence != "once", validation.Required)),
		validation.Field(&r.Recurrence, validation.Required, validation.In("once", "seconds", "minutes", "hour", "daily", "weekly"), validation.Length(1, 50)),
	)
}

type ReminderResponse struct {
	ID              uint       `json:"id"`
	RequestID       uint       `json:"request_id"`
	Message         string     `json:"message"`
	TriggeredTime   time.Time  `json:"triggered_time"`
	NextTriggerTime time.Time  `json:"next_trigger_time"`
	Occurrence      uint       `json:"occurrence"`
	Recurrence      string     `json:"recurrence"`
	Upto            *time.Time `json:"upto,omitempty"`
	CreatedAt       string     `json:"created_at"`
	UpdatedAt       string     `json:"updated_at"`
}

func (r *ReminderRequest) ToModel() (*models.Reminder, error) {
	err := r.Validate()
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrInvalidRequestBody, err)
	}

	return &models.Reminder{
		RequestID:       r.RequestID,
		Message:         r.Message,
		TriggeredTime:   r.TriggeredTime,
		Occurrence:      r.Occurrence,
		Recurrence:      r.Recurrence,
		Upto:            r.Upto,
		AfterEvery:      r.AfterEvery,
		NextTriggerTime: r.TriggeredTime,
	}, nil
}

func FromReminderModel(model *models.Reminder) *ReminderResponse {
	return &ReminderResponse{
		ID:              model.ID,
		RequestID:       model.RequestID,
		Message:         model.Message,
		TriggeredTime:   model.TriggeredTime,
		NextTriggerTime: model.NextTriggerTime,
		Occurrence:      model.Occurrence,
		Recurrence:      model.Recurrence,
		Upto:            model.Upto,
		CreatedAt:       model.CreatedAt.Format(ResponseDateFormat),
		UpdatedAt:       model.UpdatedAt.Format(ResponseDateFormat),
	}
}
