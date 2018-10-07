package pkg

import (
	"github.com/jasonlvhit/gocron"
	"gopkg.in/telegram-bot-api.v4"
	"reminder-bot/db"
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

	db.RegisterJob(chatId, remindTime)

	job := gocron.Every(1).Days().At(remindTime)
	job.Do(sendRemindMessage, reminder.BotApi, chatId)

	jobs := ReminderJobs{DailyReminder: job}
	reminder.UserJobs[chatId] = jobs
}

//Cancels the Cron-Job and removes registered Reminder-Jobs from the Reminder-Object
func (reminder Reminder) unregister(chatId int64) {
	reminderJobs, present := reminder.UserJobs[chatId]
	if present {
		gocron.Remove(reminderJobs.DailyReminder)
		delete(reminder.UserJobs, chatId)
	}
}

//TRUE - if a active Reminder-Job exists
func (reminder Reminder) hasActiveReminder(chatId int64) bool {
	_, present := reminder.UserJobs[chatId]
	return present
}

//Pauses a Reminder-Job for 7 days
func (reminder Reminder) pause(chatId int64) {
	db.UpdateSkipReminds(chatId, 7)
}

//Sends a default Message to a given chatId
func sendRemindMessage(botApi *tgbotapi.BotAPI, chatId int64) {
	skipReminds := db.GetSkipReminds(chatId)
	if skipReminds > 0 {
		db.UpdateSkipReminds(chatId, skipReminds-1)
	} else {
		msg := tgbotapi.NewMessage(chatId, "")
		msg.Text = "denk dran!"
		botApi.Send(msg)
	}

}
