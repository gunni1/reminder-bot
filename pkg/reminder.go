package pkg

import (
	"github.com/jasonlvhit/gocron"
	"gopkg.in/telegram-bot-api.v4"
)

type Reminder struct {
	BotApi   *tgbotapi.BotAPI
	UserJobs map[int64]ReminderJobs
}

type ReminderJobs struct {
	DailyReminder    *gocron.Job
	HourlyReReminder *gocron.Job
}

//Register a chatId for Reminder-Requests
func (reminder Reminder) register(chatId int64, remindTime string) {
	reminder.unregister(chatId)

	job := gocron.Every(1).Days().At(remindTime)
	job.Do(sendFirstRemindMessage, reminder.BotApi, chatId)

	jobs := ReminderJobs{DailyReminder: job}
	reminder.UserJobs[chatId] = jobs
	//gocron.Every(15).Seconds().Do(sendFirstRemindMessage, reminder.BotApi, chatId)
}

//Cancels the Cron-Job and removes registered Reminder-Jobs from the Reminder-Object
func (reminder Reminder) unregister(chatId int64) {
	reminderJobs, present := reminder.UserJobs[chatId]
	if present {
		gocron.Remove(reminderJobs.DailyReminder)
		delete(reminder.UserJobs, chatId)
	}
}

func sendFirstRemindMessage(botApi *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "")
	msg.Text = "denk dran!"
	botApi.Send(msg)
}
