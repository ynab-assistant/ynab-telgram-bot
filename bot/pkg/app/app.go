package app

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oneils/ynab-helper/bot/pkg/config"
	"github.com/oneils/ynab-helper/bot/pkg/telegram"
)

// Run starts the application
func Run(configPath string) {
	log := log.New(os.Stdout, "BOT : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal("cant init configuration for the app")
	}

	tgbotapi, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("error while creating Telegram Bot API: ", err)
	}

	tgbotapi.Debug = true

	bot := telegram.NewBot(tgbotapi, log)

	if err := bot.Start(); err != nil {
		log.Fatalf(" error while starting the bot: %v", err)
	}
}
