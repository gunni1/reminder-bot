package pkg

import (
	"github.com/jasonlvhit/gocron"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strconv"
)

var (
	GlobalUserJobs = make(map[int64]ReminderJob)
)

type Reminder struct {
	BotApi *tgbotapi.BotAPI
}

type ReminderJob struct {
	DailyReminder *gocron.Job
	SkipReminds   int
}

//Register a chatId for Reminder-Requests
func (reminder Reminder) register(chatId int64, remindTime string) {
	reminder.unregister(chatId)

	job := gocron.Every(1).Days().At(remindTime)
	job.Do(sendRemindMessage, reminder.BotApi, chatId)

	jobs := ReminderJob{DailyReminder: job, SkipReminds: 0}
	GlobalUserJobs[chatId] = jobs
}

//Cancels the Cron-Job and removes registered Reminder-Jobs from the Reminder-Object
func (reminder Reminder) unregister(chatId int64) {
	reminderJobs, present := GlobalUserJobs[chatId]
	if present {
		gocron.Remove(reminderJobs.DailyReminder)
		delete(GlobalUserJobs, chatId)
	}
}

//TRUE - if a active Reminder-Job exists
func (reminder Reminder) hasActiveReminder(chatId int64) bool {
	_, present := GlobalUserJobs[chatId]
	return present
}

//Pauses a Reminder-Job for 7 days
func (reminder Reminder) pause(chatId int64) {
	job := GlobalUserJobs[chatId]
	job.SkipReminds = 7
	GlobalUserJobs[chatId] = job
}

//Sends a default Message to a given chatId
func sendRemindMessage(botApi *tgbotapi.BotAPI, chatId int64) {
	if GlobalUserJobs[chatId].SkipReminds > 0 {
		job := GlobalUserJobs[chatId]
		job.SkipReminds = job.SkipReminds - 1
		GlobalUserJobs[chatId] = job
		log.Println("remind skipped. jetz noch: " + strconv.Itoa(GlobalUserJobs[chatId].SkipReminds))
	} else {
		msg := tgbotapi.NewMessage(chatId, "")
		msg.Text = "denk dran!"
		botApi.Send(msg)
	}

}
