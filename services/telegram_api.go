package services

import (
	"NotificationManagement/config"
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBotImpl struct {
	BotAPI       *tgbotapi.BotAPI
	active       bool
	telegramRepo domain.TelegramRepository
}

func NewTelegramAPI(repo domain.TelegramRepository) domain.TelegramAPI {
	if *config.Telegram().Enabled {
		bot, err := tgbotapi.NewBotAPI(config.Telegram().Token)
		if err == nil {
			return &TelegramBotImpl{
				active:       true,
				BotAPI:       bot,
				telegramRepo: repo,
			}
		}
	}
	return &TelegramBotImpl{
		BotAPI:       nil,
		active:       false,
		telegramRepo: repo,
	}
}

func (t *TelegramBotImpl) Start() {
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

func (t *TelegramBotImpl) HandleUpdate(update tgbotapi.Update) {
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

func (t *TelegramBotImpl) HandleMessage(message *tgbotapi.Message) {
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

func (t *TelegramBotImpl) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
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

func (t *TelegramBotImpl) SendStartMessage(msg *tgbotapi.Message) {
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

func (t *TelegramBotImpl) SendHelpMessage(chatID int64) {
	logger.Debug("Sending help message.", "chat_id", chatID)
	msg := tgbotapi.NewMessage(chatID, "This is the help message. You can use /start to see options.")
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		logger.Error("Failed to send help message", "error", err)
	}
}

func (t *TelegramBotImpl) SendMessage(chatID int64, text string, markup interface{}) {
	logger.Debug("Sending message ", "chat_id", chatID, "msg")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup

	_, err := t.BotAPI.Send(msg)
	if err != nil {
		logger.Error("Failed to send start message with keyboard", "error", err)
	}
}

func (t *TelegramBotImpl) getOTPKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Paste OTP Here", config.App().Domain),
		),
	)
	return keyboard
}
