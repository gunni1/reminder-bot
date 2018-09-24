package main

import (
	"flag"
	"github.com/jasonlvhit/gocron"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"reminder-bot/pkg"
)

func main() {
	token := flag.String("token", "x", "Telegram Bot Token")
	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	gocron.Start()

	reminder := pkg.Reminder{BotApi: bot, UserJobs: make(map[int64]pkg.ReminderJobs)}
	commandHandler := pkg.CommandHandler{BotApi: bot, Reminder: reminder}
	commandHandler.ListenForUpdates()

}
