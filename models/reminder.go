package models

import (
	"NotificationManagement/utils/datetime"
	"time"

	"gorm.io/gorm"
)

type Reminder struct {
	gorm.Model
	RequestID       uint
	Request         *CurlRequest `gorm:"foreignKey:RequestID"`
	Message         string       `gorm:"type:text;not null"`
	TriggeredTime   time.Time    `gorm:"index"`
	NextTriggerTime time.Time    `gorm:"index"`
	Occurrence      uint         `gorm:"type:int;default:0"`
	Recurrence      string       `gorm:"size:50;not null;"`
	AfterEvery      uint         `gorm:"type:int;not null"`
	TaskID          string       `gorm:"type:text"`
	Upto            *time.Time   `gorm:"index"`
}

func (r *Reminder) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*Reminder); ok {
		copyFields(r, src)
	}
}

func (r *Reminder) GetRecurrenceDuration() time.Duration {
	return datetime.RecurrenceDuration(r.AfterEvery, r.Recurrence, &r.NextTriggerTime)
}
