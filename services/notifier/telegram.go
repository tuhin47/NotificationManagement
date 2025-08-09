package notifier

import (
	"NotificationManagement/config"
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"NotificationManagement/utils"
	"NotificationManagement/utils/errutil"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strconv"
	"strings"
)

type TelegramNotifier struct {
	*tgbotapi.BotAPI
	active       bool
	telegramRepo domain.TelegramRepository
}

func NewTelegramNotifier(repo domain.TelegramRepository) domain.TelegramNotifier {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		logger.Error("Failed to create Telegram bot API", "error", err)
		return &TelegramNotifier{
			active: false,
			BotAPI: nil,
		}
	}
	t := &TelegramNotifier{
		active:       true,
		BotAPI:       bot,
		telegramRepo: repo,
	}
	return t
}

func (t *TelegramNotifier) Send(n *types.Notification) error {
	if !t.active {
		return fmt.Errorf("telegram notifier is not active")
	}
	chatID, err := strconv.ParseInt(n.To, 10, 64)
	if err != nil {
		logger.Error("Failed to parse chat ID from notification", "error", err)
		return err
	}
	msg := tgbotapi.NewMessage(chatID, n.Message)
	_, err = t.BotAPI.Send(msg) // Changed := to =
	if err != nil {
		logger.Error("Failed to send Telegram message", "error", err)
		return err
	}
	logger.Info(fmt.Sprintf("[Telegram] To: %s, Message: %s", n.To, n.Message), n)
	return nil
}

func (t *TelegramNotifier) Type() string {
	return "telegram"
}

func (t *TelegramNotifier) IsActive() bool {
	return t.active
}

func (t *TelegramNotifier) Start() {
	logger.Info("Initializing and starting Telegram bot to listen for updates.")
	if !t.active {
		logger.Info("Telegram bot is not active, skipping start.")
		return
	}

	logger.Info("Starting Telegram bot...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.BotAPI.GetUpdatesChan(u)

	for update := range updates {
		t.HandleUpdate(update)
	}
}

func (t *TelegramNotifier) VerifyOTP(ctx context.Context, otp string, userID uint) (*models.Telegram, error) {
	telegramModel, err := t.telegramRepo.GetByOTP(ctx, otp)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrRecordNotFound, fmt.Errorf("telegram chat ID not found"))
	}

	if strings.EqualFold(telegramModel.Otp, otp) {
		telegramModel.UserID = &userID
		err = t.telegramRepo.Update(ctx, telegramModel)
		if err != nil {
			return nil, errutil.NewAppError(errutil.ErrDatabaseQuery, fmt.Errorf("failed to update telegram user ID: %w", err))
		}
		return telegramModel, nil
	}
	return nil, errutil.NewAppError(errutil.ErrInvalidRequestBody, fmt.Errorf("invalid OTP"))
}

func (t *TelegramNotifier) HandleUpdate(update tgbotapi.Update) {
	logger.Debug("Processing incoming Telegram update.")
	if update.Message != nil {
		logger.Debug("Handling incoming message.")
		t.HandleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		logger.Debug("Handling callback query.")
		t.HandleCallbackQuery(update.CallbackQuery)
	}
	logger.Debug("Consider adding more handlers for other update types as needed.")
}

func (t *TelegramNotifier) HandleMessage(message *tgbotapi.Message) {
	logger.Info(fmt.Sprintf("Received message from %s: %s", message.From.UserName, message.Text))

	logger.Debug("Attempting to echo back the message.")
	msg := tgbotapi.NewMessage(message.Chat.ID, "You said: "+message.Text)
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		logger.Error("Failed to send echo message", "error", err)
	}

	logger.Debug("Registering a simple command handler.")
	switch message.Command() {
	case "start":
		t.SendStartMessage(message)
	case "help":
		t.SendHelpMessage(message.Chat.ID)
	default:
		logger.Debug("Default response for non-command messages. More complex logic can be added here.")
	}
}

func (t *TelegramNotifier) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	logger.Info(fmt.Sprintf("Received callback query from %s: %s", callbackQuery.From.UserName, callbackQuery.Data))

	logger.Debug("Responding to the callback query to remove the 'loading' state from the button.")
	callback := tgbotapi.NewCallback(callbackQuery.ID, callbackQuery.Data)
	if _, err := t.BotAPI.Request(callback); err != nil {
		logger.Error("Failed to acknowledge callback query", "error", err)
	}

	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, fmt.Sprintf("You chose: %s", callbackQuery.Data))
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		logger.Error("Failed to send message for callback query", "error", err)
	}
}

func (t *TelegramNotifier) SendStartMessage(msg *tgbotapi.Message) {
	background := context.Background()
	chatId := msg.Chat.ID
	otp := utils.GenerateRandomNumber(6)
	telegramModel := models.Telegram{
		ChatID: chatId,
		Otp:    otp,
	}
	err := t.telegramRepo.Create(background, &telegramModel)
	if err != nil {
		logger.Error("error occurred during saving", "telegram", telegramModel, err)
	}
	t.SendMessage(chatId, "Welcome to our Bot for NMS. Your OTP is: "+otp, t.getOTPKeyboard())
}

func (t *TelegramNotifier) SendHelpMessage(chatID int64) {
	logger.Debug("Sending help message.", "chat_id", chatID)
	msg := tgbotapi.NewMessage(chatID, "This is the help message. You can use /start to see options.")
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		logger.Error("Failed to send help message", "error", err)
	}
}

func (t *TelegramNotifier) SendMessage(chatID int64, text string, markup interface{}) {
	logger.Debug("Sending message ", "chat_id", chatID, "msg")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup

	_, err := t.BotAPI.Send(msg)
	if err != nil {
		logger.Error("Failed to send start message with keyboard", "error", err)
	}
}

func (t *TelegramNotifier) getOTPKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Paste OTP Here", config.App().Domain),
		),
	)
	return keyboard
}
