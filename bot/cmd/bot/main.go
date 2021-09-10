package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oneils/ynab-helper/bot/pkg/telegram"
)

func main() {

	log := log.New(os.Stdout, "BOT : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	tgbotapi, err := tgbotapi.NewBotAPI("TOKEN")
	if err != nil {
		log.Fatal("error while creating Telegram Bot API: ", err)
	}

	bot := telegram.NewBot(tgbotapi, log)

	if err := bot.Start(); err != nil {
		log.Fatalf(" error while starting the bot: %v", err)
	}
}
