package pkg

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type CommandHandler struct {
	BotApi   *tgbotapi.BotAPI
	Reminder Reminder
}

func (handler CommandHandler) ListenForUpdates() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, _ := handler.BotApi.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "/remindme 08:00 -> folgenden 3 Wochen jeden Tag um 8 Uhr eine Erinnerung."
			case "remindme":
				handler.Reminder.register(update.Message.Chat.ID)
				msg.Text = "you will be reminded"
			case "stop":
				handler.Reminder.unregister(update.Message.Chat.ID)
			case "ok":
				msg.Text = "I'm ok."
			default:
				msg.Text = "I don't know that command"
			}
			handler.BotApi.Send(msg)
		}

	}
}
