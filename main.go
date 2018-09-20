package main

import (
	"flag"
	"github.com/robfig/cron"
	"gopkg.in/telegram-bot-api.v4"
	"log"
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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	c := cron.New()
	c.Start()

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

				register(bot, c, update.Message.Chat.ID)
			case "ok":
				msg.Text = "I'm ok."
			default:
				msg.Text = "I don't know that command"
			}
			//bot.Send(msg)
		}

	}
}

func register(bot *tgbotapi.BotAPI, cron *cron.Cron, chatId int64) {
	cron.AddFunc("30 * * * * * ", func() {
		msg := tgbotapi.NewMessage(chatId, "aa")
		msg.Text = "tick"
		bot.Send(msg)
	})

}
