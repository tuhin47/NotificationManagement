package models

import (
	"time"

	"gorm.io/gorm"
)

type Reminder struct {
	gorm.Model
	RequestID       uint      `gorm:"index"` // Foreign key to CurlRequest
	Message         string    `gorm:"type:text;not null"`
	TriggeredTime   time.Time `gorm:"index"`
	NextTriggerTime time.Time `gorm:"index"`
	Occurrence      uint      `gorm:"type:int;default:0"`
	Recurrence      string    `gorm:"size:50;check:recurrence IN ('once','minutes','hour','daily','weekly')"`
}

func (r *Reminder) UpdateFromModel(source ModelInterface) {
	if src, ok := source.(*Reminder); ok {
		copyFields(r, src)
	}
}
