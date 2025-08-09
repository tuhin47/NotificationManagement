package worker

import (
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type ReminderTaskHandler struct {
	reminderService        domain.ReminderService
	asynqService           domain.AsynqService
	notificationDispatcher domain.NotificationDispatcher
}

func NewReminderTaskHandler(reminderService domain.ReminderService, asynqService domain.AsynqService, dispatcher domain.NotificationDispatcher) *ReminderTaskHandler {

	return &ReminderTaskHandler{
		reminderService:        reminderService,
		asynqService:           asynqService,
		notificationDispatcher: dispatcher,
	}
}

func (h *ReminderTaskHandler) HandleReminderTask(ctx context.Context, task *asynq.Task) error {
	var reminder models.Reminder
	if err := json.Unmarshal(task.Payload(), &reminder); err != nil {
		logger.Error("Failed to unmarshal reminder payload", "error", err)
		return fmt.Errorf("failed to unmarshal reminder payload: %w", err)
	}

	logger.Info("Processing reminder task", "reminder_id", reminder.ID, "message", reminder.Message)
	err := h.reminderService.SendReminders(ctx, reminder.ID)

	if err != nil {
		logger.Error("Reminder Has issues", err)
	}

	if reminder.Recurrence != "once" {
		nextTrigger := reminder.NextTriggerTime.Add(reminder.GetRecurrenceDuration())

		if reminder.Upto != nil && nextTrigger.After(*reminder.Upto) {
			logger.Info("Reminder reached 'Upto' time, deleting", "reminder_id", reminder.ID)
			err := h.reminderService.DeleteModel(ctx, reminder.ID)
			if err != nil {
				logger.Error("Failed to delete reminder after reaching 'Upto' time", "error", err, "reminder_id", reminder.ID)
				return err
			}
			logger.Info("Reminder deleted after reaching 'Upto' time", "reminder_id", reminder.ID)
			return nil
		}

		reminder.TriggeredTime = reminder.NextTriggerTime
		reminder.NextTriggerTime = nextTrigger
		reminder.Occurrence++

		_, err := h.reminderService.UpdateModel(ctx, reminder.ID, &reminder)
		if err != nil {
			logger.Error("Failed to update recurring reminder", "error", err, "reminder_id", reminder.ID)
			return fmt.Errorf("failed to update reminder: %w", err)
		}
		logger.Info("Recurring reminder updated and scheduled for next occurrence", "reminder_id", reminder.ID, "next_trigger_time", nextTrigger)

		_, err = h.asynqService.CreateReminderTask(ctx, &reminder)
		if err != nil {
			logger.Error("Failed to schedule next reminder task", "error", err, "reminder_id", reminder.ID)
			return fmt.Errorf("failed to schedule next reminder task: %w", err)
		}
	} else {
		reminder.TriggeredTime = reminder.NextTriggerTime
		reminder.Occurrence++

		_, err := h.reminderService.UpdateModel(ctx, reminder.ID, &reminder)
		if err != nil {
			logger.Error("Failed to update one-time reminder", "error", err, "reminder_id", reminder.ID)
			return fmt.Errorf("failed to update reminder: %w", err)
		}
		logger.Info("One-time reminder marked as triggered", "reminder_id", reminder.ID)
	}

	return nil
}
