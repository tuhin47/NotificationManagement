package server

import (
	"NotificationManagement/conn"
	"NotificationManagement/controllers"
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/middleware"
	"NotificationManagement/repositories"
	"NotificationManagement/routes"
	"NotificationManagement/services"
	"NotificationManagement/services/notifier"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"gorm.io/gorm" // Import gorm
)

func NewEcho(db *gorm.DB) *echo.Echo {
	e := echo.New()
	e.Use(interceptLogger)
	//e.Use(errutil.ErrorHandler())
	e.Use(middleware.TransactionMiddleware(db))
	return e
}

func RegisterRoutes(e *echo.Echo, curlController domain.CurlController, llmController domain.LLMController, reminderController domain.ReminderController, aiController domain.AIRequestController, userController domain.UserController, notificationController *controllers.NotificationController, userService domain.UserService, telegramController domain.TelegramController) {
	keycloakMiddleware := middleware.KeycloakMiddleware(userService)
	routes.RegisterCurlRoutes(e, curlController, &keycloakMiddleware)
	routes.RegisterLLMRoutes(e, llmController, &keycloakMiddleware)
	routes.RegisterReminderRoutes(e, reminderController, &keycloakMiddleware)
	routes.RegisterAIRoutes(e, aiController, &keycloakMiddleware)
	routes.RegisterUserRoutes(e, userController, &keycloakMiddleware)
	routes.RegisterTelegramRoutes(e, telegramController, &keycloakMiddleware)
	routes.RegisterNotificationRoutes(e, notificationController, &keycloakMiddleware)
}

func interceptLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		logger.Info("Incoming request",
			"method", req.Method,
			"path", req.URL.Path,
			"query", req.URL.RawQuery,
			"remote_ip", c.RealIP(),
			"user_agent", req.UserAgent(),
		)
		return next(c)
	}
}

var Module = fx.Options(
	fx.Provide(
		NewEcho,
		conn.NewDB,
		conn.NewAsynq,
		conn.NewAsynqInspector,
		notifier.NewEmailNotifier,
		notifier.NewSMSNotifier,
		notifier.NewTelegramNotifier,
		notifier.NewNotificationDispatcher,

		controllers.NewAIRequestController,
		controllers.NewCurlController,
		controllers.NewLLMController,
		controllers.NewReminderController,
		controllers.NewUserController,
		controllers.NewNotificationController,
		controllers.NewTelegramController,

		repositories.NewAIModelRepository,
		repositories.NewCurlRequestRepository,
		repositories.NewAdditionalFieldsRepository,
		repositories.NewDeepseekModelRepository,
		repositories.NewGeminiRepository,
		repositories.NewLLMRepository,
		repositories.NewReminderRepository,
		repositories.NewUserRepository,
		repositories.NewTelegramRepository,

		services.NewAIModelService,
		services.NewAsynqService,
		services.NewCurlService,
		services.NewDeepseekModelService,
		services.NewGeminiService,
		services.NewLLMService,
		services.NewReminderService,
		services.NewUserService,
		services.NewAIDispatcher,
		services.NewTelegramAPI, // TODO : Need to remove from here.Manage By Worker
	),
	fx.Invoke(RegisterRoutes),
)
