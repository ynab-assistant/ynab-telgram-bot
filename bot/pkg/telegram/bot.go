package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oneils/ynab-helper/bot/pkg/transaction"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	log *log.Logger
	txn *transaction.Transaction
}

func NewBot(bot *tgbotapi.BotAPI, log *log.Logger, txn *transaction.Transaction) *Bot {
	return &Bot{
		bot: bot,
		log: log,
		txn: txn,
	}
}

func (b *Bot) Start() error {

	updates := b.initUpdatesChannel()

	b.handleUdate(updates)

	return nil
}

func (b *Bot) handleUdate(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.log.Printf("error while handling a command: %s. Error: %v", update.Message.Command(), err)
			}
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	b.log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
