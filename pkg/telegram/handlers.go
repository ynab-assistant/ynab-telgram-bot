package telegram

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oneils/ynab-helper/bot/pkg/transaction"
)

const commandStart = "start"

func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	txnMsg := transaction.TxnMessage{
		ChatID:   message.Chat.ID,
		UserName: message.From.UserName,
		Text:     message.Text,
	}

	err := b.txn.Save(context.Background(), txnMsg, time.Now().UTC())
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
		msg.ReplyToMessageID = message.MessageID
		_, err = b.bot.Send(msg)
		return err
	}

	b.logger.Printf("Verification: [%s] %d", message.From.UserName, message.Chat.ID)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Saved")

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Start instructions will be here later")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Start instructions will be here later")
	_, err := b.bot.Send(msg)
	return err
}
