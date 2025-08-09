package conn

import (
	"NotificationManagement/config"
	"github.com/hibiken/asynq"
)

func NewAsynq() *asynq.Client {
	asynqConfig := config.Asynq()
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     asynqConfig.RedisAddr,
		DB:       *asynqConfig.DB,
		Password: asynqConfig.Pass,
	})
}

func NewAsynqInspector() *asynq.Inspector {
	asynqConfig := config.Asynq()
	return asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     asynqConfig.RedisAddr,
		DB:       *asynqConfig.DB,
		Password: asynqConfig.Pass,
	})
}
