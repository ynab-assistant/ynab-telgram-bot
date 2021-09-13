package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oneils/ynab-helper/bot/pkg/transaction"
)

// Bot contains required dependencies for running Telegram Bot
type Bot struct {
	bot    *tgbotapi.BotAPI
	logger *log.Logger
	txn    *transaction.Transaction
}

// NewBot a helper for creating a Bot
func NewBot(bot *tgbotapi.BotAPI, logger *log.Logger, txn *transaction.Transaction) *Bot {
	return &Bot{
		bot:    bot,
		logger: logger,
		txn:    txn,
	}
}

// Start runs the Bot and handles updates from the Bot
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
				b.logger.Printf("error while handling a command: %s. Error: %v", update.Message.Command(), err)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.logger.Printf("cant handle Telegram API message. Error: %v", err)
		}
	}
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	b.logger.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
