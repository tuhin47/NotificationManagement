package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"context"
)

type ReminderServiceImpl struct {
	Repo domain.ReminderRepository
}

func NewReminderService(repo domain.ReminderRepository) domain.ReminderService {
	return &ReminderServiceImpl{Repo: repo}
}

func (s *ReminderServiceImpl) CreateReminder(reminder *models.Reminder) error {
	return s.Repo.Create(context.Background(), reminder)
}

func (s *ReminderServiceImpl) GetReminderByID(id uint) (*models.Reminder, error) {
	reminder, err := s.Repo.GetByID(context.Background(), id, nil)
	if err != nil {
		return nil, err
	}
	return reminder, nil
}

func (s *ReminderServiceImpl) GetAllReminders(limit, offset int) ([]models.Reminder, error) {
	reminders, err := s.Repo.GetAll(context.Background(), limit, offset)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

func (s *ReminderServiceImpl) UpdateReminder(id uint, reminder *models.Reminder) error {
	// First check if the record exists
	existing, err := s.Repo.GetByID(context.Background(), id, nil)
	if err != nil {
		return err
	}

	// Update the existing record with new data
	existing.RequestID = reminder.RequestID
	existing.Message = reminder.Message
	existing.TriggeredTime = reminder.TriggeredTime
	existing.NextTriggerTime = reminder.NextTriggerTime
	existing.Occurrence = reminder.Occurrence
	existing.Recurrence = reminder.Recurrence

	// Save the updated record
	err = s.Repo.Update(context.Background(), existing)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReminderServiceImpl) DeleteReminder(id uint) error {
	err := s.Repo.Delete(context.Background(), id)
	if err != nil {
		return err
	}

	return nil
}
