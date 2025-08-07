package domain

import (
	"NotificationManagement/models"
	"NotificationManagement/types"
	"context"
	"github.com/hibiken/asynq"
	"time"
)

type (
	AsynqRepository interface {
		CreateTask(event types.AsynqTaskType, payload interface{}) (*asynq.Task, error)
		EnqueueTask(task *asynq.Task, customOpts *types.AsynqOption) (string, error)
		DequeueTask(taskID string) error
	}

	AsynqService interface {
		CreateReminderTask(ctx context.Context, reminder *models.Reminder) (string, error)
		UpdateReminderTask(ctx context.Context, reminder *models.Reminder) error
		CancelReminderTask(ctx context.Context, reminderID uint) error
		GetTaskInfo(ctx context.Context, taskID string) (interface{}, error)
		ScheduleTask(ctx context.Context, taskType string, payload interface{}, processAt time.Time, opts ...interface{}) (string, error)
	}
)
