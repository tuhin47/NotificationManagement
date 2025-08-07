package types

import (
	"time"
)

type (
	AsynqOption struct {
		TaskID           string
		Retry            int
		Queue            string
		RetentionHours   time.Duration
		DelaySeconds     time.Duration
		UniqueTTLSeconds time.Duration
	}

	AsynqTaskType string
)

func (t AsynqTaskType) String() string {
	return string(t)
}

const (
	AsynqTaskTypeHandleReminder AsynqTaskType = "go:nms:reminder"
)
