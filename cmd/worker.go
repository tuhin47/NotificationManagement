package cmd

import (
	"NotificationManagement/config"
	"NotificationManagement/conn"
	"NotificationManagement/domain"
	"NotificationManagement/repositories"
	"NotificationManagement/services"
	"NotificationManagement/services/notifier"
	"NotificationManagement/types"
	"NotificationManagement/worker"
	"context"
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"log"
	"time"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the asynq task worker",
	Run: func(cmd *cobra.Command, args []string) {
		app := fx.New(
			fx.Provide(
				NewAsynqServer,
				conn.NewDB,
				conn.NewAsynq,
				conn.NewAsynqInspector,

				repositories.NewReminderRepository,
				repositories.NewAIModelRepository,
				repositories.NewTelegramRepository,
				repositories.NewUserRepository,
				repositories.NewCurlRequestRepository,
				repositories.NewGeminiRepository,
				repositories.NewDeepseekModelRepository,
				repositories.NewAdditionalFieldsRepository,

				services.NewReminderService,
				services.NewAsynqService,
				services.NewUserService,
				services.NewTelegramAPI,
				services.NewGeminiService,
				services.NewDeepseekModelService,
				services.NewAIDispatcher,
				services.NewCurlService,
				services.NewAIModelService,

				worker.NewReminderTaskHandler,

				notifier.NewEmailNotifier,
				notifier.NewSMSNotifier,
				notifier.NewTelegramNotifier,
				notifier.NewNotificationDispatcher,
			),
			fx.Invoke(registerWorker),
			fx.Invoke(registerHooks),
		)

		startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.Start(startCtx); err != nil {
			log.Fatalf("failed to start fx app: %v", err)
		}

		<-app.Done()

		stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.Stop(stopCtx); err != nil {
			log.Fatalf("failed to stop fx app: %v", err)
		}
	},
}

func registerHooks(telegramAPI domain.TelegramAPI) {
	go telegramAPI.Start()
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
					return time.Duration(*config.Asynq().EventReminderTaskRetryDelay) * time.Second
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
