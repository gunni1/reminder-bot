package main

import (
	"flag"
	"github.com/jasonlvhit/gocron"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"reminder-bot/db"
	"reminder-bot/pkg"
)

func main() {
	token := flag.String("token", "x", "Telegram Bot Token")
	dbUrl := flag.String("mongodb", "localhost:27017", "Connection String to a MongoDB")
	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	session := db.OpenDbConnection(*dbUrl)
	defer session.Close()

	gocron.Start()

	//Aus DB initialisieren?
	reminderJobsPerUser := make(map[int64]pkg.ReminderJobs)

	reminder := pkg.Reminder{BotApi: bot, UserJobs: reminderJobsPerUser}
	commandHandler := pkg.CommandHandler{BotApi: bot, Reminder: reminder}
	commandHandler.ListenForUpdates()

}
