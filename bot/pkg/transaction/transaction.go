package transaction

import (
	"context"
	"log"
	"time"

	"github.com/oneils/ynab-helper/bot/pkg/sms"
	"github.com/oneils/ynab-helper/bot/pkg/transaction/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Transaction struct {
	DB                   *mongo.Database
	log                  *log.Logger
	parser               sms.Parser
	txnRepository        repository.TxnRepository
	invalidSmsRepository repository.InvalidSmsRepository
}

func New(log *log.Logger, parser sms.Parser, txnRepository repository.TxnRepository, invalidSmsRepo repository.InvalidSmsRepository) *Transaction {
	return &Transaction{
		log:                  log,
		parser:               parser,
		txnRepository:        txnRepository,
		invalidSmsRepository: invalidSmsRepo,
	}
}

func (t *Transaction) Save(ctx context.Context, txnMsg TxnMessage, now time.Time) error {
	msg, err := t.parser.Parse(txnMsg.Text)
	if err != nil {
		log.Printf("[ERROR] can't parse incomming sms message. Transaction was not saved to DB. \n\t\tSMS: %s\n\t\tSMSEror: %v", txnMsg.Text, err)

		invalidSms := repository.InvalidSmsRecord{
			ChatID:      txnMsg.ChatID,
			UserName:    txnMsg.UserName,
			SmsMessage:  txnMsg.Text,
			DateCreated: time.Now(),
			DateUpdated: time.Now(),
		}

		if err := t.invalidSmsRepository.Save(ctx, invalidSms); err != nil {
			return err
		}

		return err
	}

	txnDate := time.Now() // parse from msg.Transaction.Date
	newTxn := repository.NewTXNRecord{
		ChatID:      txnMsg.ChatID,
		UserName:    txnMsg.UserName,
		BankName:    msg.BankName,
		CardNumber:  msg.CardNumber,
		TxnDate:     txnDate,
		Type:        msg.Transaction.Type,
		Currency:    msg.Currency,
		Amount:      msg.Amount,
		Payee:       msg.Payee,
		CountryCode: msg.CountryCode,
		SmsMessage:  msg.OriginalMsg,
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}
	err = t.txnRepository.Save(ctx, newTxn)
	if err != nil {
		return errors.Wrap(err, "can't save transaction to DB")
	}

	return nil
}
