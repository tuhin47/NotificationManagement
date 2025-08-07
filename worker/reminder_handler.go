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

// ReminderTaskHandler handles reminder notification tasks
type ReminderTaskHandler struct {
	reminderService domain.ReminderService
	asynqService    domain.AsynqService
}

// NewReminderTaskHandler creates a new handler for reminder tasks
func NewReminderTaskHandler(reminderService domain.ReminderService, asynqService domain.AsynqService) *ReminderTaskHandler {
	return &ReminderTaskHandler{
		reminderService: reminderService,
		asynqService:    asynqService,
	}
}

func (h *ReminderTaskHandler) HandleReminderTask(ctx context.Context, task *asynq.Task) error {
	var reminder models.Reminder
	if err := json.Unmarshal(task.Payload(), &reminder); err != nil {
		return fmt.Errorf("failed to unmarshal reminder payload: %w", err)
	}

	logger.Info("Processing reminder task", "reminder_id", reminder.ID, "message", reminder.Message)

	// Check if this is a recurring reminder
	if reminder.Recurrence != "once" {
		// Calculate next trigger time
		nextTrigger := reminder.NextTriggerTime.Add(reminder.GetRecurrenceDuration())

		// Check if the next trigger time is beyond the 'Upto' time
		if reminder.Upto != nil && nextTrigger.After(*reminder.Upto) {
			logger.Info("Reminder reached 'Upto' time, deleting", "reminder_id", reminder.ID)
			err := h.reminderService.DeleteModel(ctx, reminder.ID)
			if err != nil {
				logger.Error("Failed to delete reminder after reaching 'Upto' time", "error", err, "reminder_id", reminder.ID)
				return err
			}
			return nil
		}

		// Update the reminder
		reminder.TriggeredTime = reminder.NextTriggerTime
		reminder.NextTriggerTime = nextTrigger
		reminder.Occurrence++

		// Save the updated reminder
		_, err := h.reminderService.UpdateModel(ctx, reminder.ID, &reminder)
		if err != nil {
			return fmt.Errorf("failed to update reminder: %w", err)
		}

		// Schedule the next task
		_, err = h.asynqService.CreateReminderTask(ctx, &reminder)
		if err != nil {
			return fmt.Errorf("failed to schedule next reminder task: %w", err)
		}
	} else {
		// For non-recurring reminders, mark as triggered
		reminder.TriggeredTime = reminder.NextTriggerTime
		reminder.Occurrence++

		// Save the updated reminder
		_, err := h.reminderService.UpdateModel(ctx, reminder.ID, &reminder)
		if err != nil {
			return fmt.Errorf("failed to update reminder: %w", err)
		}
	}

	return nil
}
