package conn

import (
	"NotificationManagement/config"
	"github.com/hibiken/asynq"
)

func NewAsynq() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     config.Asynq().RedisAddr,
		DB:       config.Asynq().DB,
		Password: config.Asynq().Pass,
	})
}

func NewAsynqInspector() *asynq.Inspector {
	return asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     config.Asynq().RedisAddr,
		DB:       config.Asynq().DB,
		Password: config.Asynq().Pass,
	})
}
