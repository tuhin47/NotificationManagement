package conn

import (
	"NotificationManagement/services/notifier"

	"go.uber.org/fx"
)

func ProvideNotifiers() fx.Option {
	return fx.Options(
		fx.Provide(notifier.NewEmailNotifier),
		fx.Provide(notifier.NewSMSNotifier),
		fx.Provide(notifier.NewTelegramNotifier),
		fx.Provide(notifier.NewDispatcher),
	)
}
