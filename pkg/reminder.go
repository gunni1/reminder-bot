package pkg

import (
	"github.com/robfig/cron"
	"gopkg.in/telegram-bot-api.v4"
)

type Reminder struct {
	Cron   *cron.Cron
	BotApi *tgbotapi.BotAPI
}

//Sends
func (reminder Reminder) register(chatId int64) {

	reminder.Cron.AddFunc("30 * * * * * ", func() {

		//register 1 hour reminder for today

		msg := tgbotapi.NewMessage(chatId, "")
		msg.Text = "tick"
		reminder.BotApi.Send(msg)
	})
}
func (reminder Reminder) unregister(chatId int64) {

}
