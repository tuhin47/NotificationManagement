package services

import (
	"NotificationManagement/config"
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

type AsynqServiceImpl struct {
	client    *asynq.Client
	inspector *asynq.Inspector
	repo      domain.ReminderRepository
}

func NewAsynqService(repo domain.ReminderRepository, client *asynq.Client, inspector *asynq.Inspector) domain.AsynqService {
	return &AsynqServiceImpl{
		client:    client,
		inspector: inspector,
		repo:      repo,
	}
}

func (s *AsynqServiceImpl) CreateReminderTask(ctx context.Context, reminder *models.Reminder) (string, error) {
	reminderTime := reminder.NextTriggerTime
	now := time.Now().UTC()

	if reminder.Upto.Before(now) {
		logger.Info("Current time is beyond the 'Upto' time, skipping event reminder for event: ", reminder.Message)
		_ = s.CancelReminderTask(ctx, reminder.ID)
		return "", errutil.NewAppError(errutil.ErrInvalidRequestBody, fmt.Errorf("error reminder time is beyond 'Upto' time"))
	}
	if reminderTime.Before(now) {
		logger.Warn("Reminder time is in the past, considering next occurrence or cancellation for event: ", reminder.Message)
		_ = s.CancelReminderTask(ctx, reminder.ID)
		return "", errutil.NewAppError(errutil.ErrInvalidRequestBody, fmt.Errorf("error time already passed"))
	}

	logger.Info("Enqueuing event reminder  for event: ", reminder.Message)

	logger.Debug("Converting reminder to payload")
	payload, err := json.Marshal(reminder)
	if err != nil {
		return "", errutil.NewAppError(errutil.ErrTaskMarshalPayloadFailed, err)
	}

	logger.Debug("Creating asynq task")
	task := asynq.NewTask(types.AsynqTaskTypeHandleReminder.String(), payload)
	logger.Debug("Configuring asynq task options")
	opts := []asynq.Option{
		asynq.Queue(config.Asynq().Queue),
		asynq.ProcessAt(reminderTime),
		asynq.MaxRetry(*config.Asynq().EventReminderTaskRetryCount),
		asynq.Retention(config.Asynq().Retention),
	}
	logger.Debug("Enqueuing asynq task")
	info, err := s.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return "", errutil.NewAppError(errutil.ErrTaskEnqueueFailed, err)
	}
	logger.Info("Created async task for reminder", "reminder_id", reminder.ID, "task_id", info.ID)
	return info.ID, nil
}

// UpdateReminderTask updates an existing asynq task
func (s *AsynqServiceImpl) UpdateReminderTask(ctx context.Context, reminder *models.Reminder) error {
	logger.Debug("Checking if task ID is empty to create a new task")
	if reminder.TaskID == "" {
		taskID, err := s.CreateReminderTask(ctx, reminder)
		if err != nil {
			return err
		}

		logger.Debug("Updating the reminder with the new task ID")
		reminder.TaskID = taskID
		return s.repo.Update(ctx, reminder)
	}

	logger.Debug("Task ID exists, cancelling the existing task and creating a new one")
	err := s.CancelReminderTask(ctx, reminder.ID)
	if err != nil {
		logger.Warn("Failed to cancel existing task", "error", err, "reminder_id", reminder.ID)
	}

	logger.Debug("Creating new task")
	taskID, err := s.CreateReminderTask(ctx, reminder)
	if err != nil {
		return err
	}

	logger.Debug("Updating the reminder with the new task ID")
	reminder.TaskID = taskID
	return s.repo.Update(ctx, reminder)
}

// CancelReminderTask cancels an existing reminder task
func (s *AsynqServiceImpl) CancelReminderTask(ctx context.Context, reminderID uint) error {
	logger.Debug("Getting the reminder to find the task ID", "reminder_id", reminderID)
	reminder, err := s.repo.GetByID(ctx, reminderID, nil)
	if err != nil {
		return errutil.NewAppError(errutil.ErrRecordNotFound, err)
	}

	if reminder.TaskID == "" {
		logger.Debug("No task to cancel for reminder", "reminder_id", reminderID)
		return nil
	}

	logger.Debug("Attempting to cancel and delete task", "task_id", reminder.TaskID)
	if err := s.inspector.CancelProcessing(reminder.TaskID); err != nil {
		logger.Debug("Task not processing, attempting to delete it directly", "task_id", reminder.TaskID)
		if err := s.inspector.DeleteTask(config.Asynq().Queue, reminder.TaskID); err != nil {
			return errutil.NewAppError(errutil.ErrTaskDeletionFailed, err)
		}
	} else {
		logger.Debug("Task cancellation succeeded, also attempting to delete it", "task_id", reminder.TaskID)
		if err := s.inspector.DeleteTask(config.Asynq().Queue, reminder.TaskID); err != nil {
			return errutil.NewAppError(errutil.ErrTaskDeletionFailed, err)
		}
	}

	logger.Info("Cancelled and deleted reminder task", "reminder_id", reminderID, "task_id", reminder.TaskID)

	logger.Debug("Clearing the task ID from the reminder", "reminder_id", reminderID)
	reminder.TaskID = ""
	return s.repo.Update(ctx, reminder)
}

// GetTaskInfo retrieves information about a reminder task
func (s *AsynqServiceImpl) GetTaskInfo(ctx context.Context, taskID string) (interface{}, error) {
	info, err := s.inspector.GetTaskInfo(config.Asynq().Queue, taskID)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrTaskInfoRetrievalFailed, err)
	}
	return info, nil
}

// ScheduleTask schedules a generic task with provided payload and options
func (s *AsynqServiceImpl) ScheduleTask(ctx context.Context, taskType string, payload interface{}, processAt time.Time, opts ...interface{}) (string, error) {
	logger.Debug("Converting payload to JSON for generic task", "task_type", taskType)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", errutil.NewAppError(errutil.ErrTaskMarshalPayloadFailed, err)
	}

	logger.Debug("Creating generic asynq task", "task_type", taskType)
	task := asynq.NewTask(taskType, payloadBytes)

	logger.Debug("Converting generic options to asynq options")
	var asynqOpts []asynq.Option
	asynqOpts = append(asynqOpts, asynq.Queue(config.Asynq().Queue))
	asynqOpts = append(asynqOpts, asynq.ProcessAt(processAt))

	logger.Debug("Adding additional provided options to generic task")
	for _, opt := range opts {
		if asynqOpt, ok := opt.(asynq.Option); ok {
			asynqOpts = append(asynqOpts, asynqOpt)
		}
	}

	logger.Debug("Enqueuing generic asynq task", "task_type", taskType)
	info, err := s.client.EnqueueContext(ctx, task, asynqOpts...)
	if err != nil {
		return "", errutil.NewAppError(errutil.ErrTaskEnqueueFailed, err)
	}

	return info.ID, nil
}
