package cmd

import (
	"NotificationManagement/config"
	"NotificationManagement/conn"
	"NotificationManagement/repositories"
	"NotificationManagement/services"
	"NotificationManagement/types"
	"NotificationManagement/worker"
	"context"
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"log"
	"time"
)

// Inside your worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the asynq task worker",
	Run: func(cmd *cobra.Command, args []string) {
		app := fx.New(
			// Provide config, logger, etc. if needed
			fx.Provide(
				NewAsynqServer,
				conn.NewDB,
				conn.NewAsynq,
				conn.NewAsynqInspector,
				services.NewReminderService,
				services.NewAsynqService,
				repositories.NewReminderRepository,
				worker.NewReminderTaskHandler,
			),
			fx.Invoke(registerWorker),
		)

		startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.Start(startCtx); err != nil {
			log.Fatalf("failed to start fx app: %v", err)
		}

		// Wait for interrupt signal or similar
		<-app.Done()

		stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.Stop(stopCtx); err != nil {
			log.Fatalf("failed to stop fx app: %v", err)
		}
	},
}

func NewAsynqServer() *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.Asynq().RedisAddr,
			DB:       *config.Asynq().DB,
			Password: config.Asynq().Pass,
		},
		asynq.Config{
			Concurrency: *config.Asynq().Concurrency,
			Queues: map[string]int{
				config.Asynq().Queue: 1,
			},
			RetryDelayFunc: func(numOfRetry int, e error, t *asynq.Task) time.Duration {
				switch t.Type() {
				case types.AsynqTaskTypeHandleReminder.String():
					return config.Asynq().EventReminderTaskRetryDelay * time.Second
				default:
					return asynq.DefaultRetryDelayFunc(numOfRetry, e, t)
				}
			},
		},
	)
}

func registerWorker(
	lifecycle fx.Lifecycle,
	server *asynq.Server,
	handler *worker.ReminderTaskHandler,
) {
	mux := asynq.NewServeMux()
	mux.HandleFunc(types.AsynqTaskTypeHandleReminder.String(), handler.HandleReminderTask)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Run(mux); err != nil {
					log.Printf("Asynq server exited with error: %v", err)
				}
			}()
			log.Println("Asynq worker started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping asynq worker...")
			server.Shutdown()
			return nil
		},
	})
}
