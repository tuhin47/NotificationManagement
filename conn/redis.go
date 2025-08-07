package conn

import (
	"NotificationManagement/config"
	"NotificationManagement/logger"
	"fmt"
	"github.com/go-redis/redis"
)

var client *redis.Client

func ConnectRedis() {
	conf := config.Redis()

	logger.Info("connecting to redis at ", "host", conf.Host, "port", conf.Port)

	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	if _, err := client.Ping().Result(); err != nil {
		logger.Error("failed to connect redis: ", err)
		panic(err)
	}

	logger.Info("redis connection successful...")
}

func Redis() *redis.Client {
	return client
}
