package domain

import (
	"NotificationManagement/types"
	"context"
)

type Notifier interface {
	Send(context.Context, *types.Notification) error
	Type() string
	IsActive() bool
}

type NotificationDispatcher interface {
	Notify(ctx context.Context, n *types.Notification) error
}
