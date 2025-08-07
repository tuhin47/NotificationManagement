package models

import (
	"time"

	"gorm.io/gorm"
)

type Reminder struct {
	gorm.Model
	RequestID       uint       `gorm:"index"` // Foreign key to CurlRequest
	Message         string     `gorm:"type:text;not null"`
	TriggeredTime   time.Time  `gorm:"index"`
	NextTriggerTime time.Time  `gorm:"index"`
	Occurrence      uint       `gorm:"type:int;default:0"`
	Recurrence      string     `gorm:"size:50;not null;check:recurrence IN ('once','minutes','hour','daily','weekly')"`
	AfterEvery      uint       `gorm:"type:int;not null"`
	TaskID          string     `gorm:"type:text"`
	Upto            *time.Time `gorm:"index"`
}

func (r *Reminder) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*Reminder); ok {
		copyFields(r, src)
	}
}

func (r *Reminder) GetRecurrenceDuration() time.Duration {
	var minutes uint
	switch r.Recurrence {
	case "minutes":
		minutes = r.AfterEvery
	case "hour":
		minutes = r.AfterEvery * 60
	case "daily":
		minutes = r.AfterEvery * 24 * 60
	case "weekly":
		minutes = r.AfterEvery * 7 * 24 * 60
	case "once":
		minutes = 0 // For 'once' recurrence, no minutes for next trigger
	default:
		minutes = 0 // Default to 0 or handle error
	}
	return time.Duration(minutes) * time.Minute
}
