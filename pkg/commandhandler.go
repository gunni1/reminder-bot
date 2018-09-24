package pkg

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"regexp"
)

type CommandHandler struct {
	BotApi   *tgbotapi.BotAPI
	Reminder Reminder
}

//Listens on Chat-Bot-Updates und responds to known commands
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
				msg.Text = "/remindme 08:00 -> jeden Tag um 8 Uhr eine Erinnerung."
			case "remindme":
				msg.Text = handler.handleRemindMeCommand(update.Message.CommandArguments(), update.Message.Chat.ID)
			case "stop":
				handler.Reminder.unregister(update.Message.Chat.ID)
				msg.Text = "you will not longer be reminded."
			case "ok":
				msg.Text = "I'm ok."
			default:
				msg.Text = "I don't know that command"
			}
			handler.BotApi.Send(msg)
		}

	}
}

func (handler CommandHandler) handleRemindMeCommand(argument string, chatId int64) string {
	if isValidRemindTime(argument) {
		handler.Reminder.register(chatId, argument)
		return "you will be reminded at " + argument
	} else {
		return "Invalid time format. Please use HH:MM e.G.: /remindme 08:00"
	}
}

func isValidRemindTime(remindTime string) bool {
	match, _ := regexp.MatchString("^([0-9]|0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$", remindTime)
	return match
}
