package pkg

import "github.com/robfig/cron"

var (
	c cron.Cron
)

func addReminder(chatId int64) {
	c.AddFunc("0 0 8 * * * ", func() {})

}
