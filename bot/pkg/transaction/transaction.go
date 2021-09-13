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

// Transaction has required dependencies to manipulate transactions
type Transaction struct {
	DB             *mongo.Database
	logger         *log.Logger
	parser         sms.Parser
	txnRepo        repository.TXNer
	invalidSmsRepo repository.InvalidSMSer
}

// New creates a Transaction
func New(logger *log.Logger, parser sms.Parser, txnRepository repository.TXNer, invalidSmsRepo repository.InvalidSMSer) *Transaction {
	return &Transaction{
		logger:         logger,
		parser:         parser,
		txnRepo:        txnRepository,
		invalidSmsRepo: invalidSmsRepo,
	}
}

// Save stores the Transaction to the DB
func (t *Transaction) Save(ctx context.Context, txnMsg TxnMessage, now time.Time) error {
	msg, err := t.parser.Parse(txnMsg.Text)
	if err != nil {
		t.logger.Printf("[ERROR] can't parse incomming sms message. Transaction was not saved to DB. \n\t\tSMS: %s\n\t\tSMSEror: %v", txnMsg.Text, err)

		invalidSms := repository.InvalidSmsRecord{
			ChatID:      txnMsg.ChatID,
			UserName:    txnMsg.UserName,
			SmsMessage:  txnMsg.Text,
			DateCreated: now,
			DateUpdated: now,
		}

		if saveErr := t.invalidSmsRepo.Save(ctx, &invalidSms); saveErr != nil {
			t.logger.Printf("[ERROR] can't save invalid sms to the DB. \n\t\tinvalidSMS: %v\n\t\tError: %v", invalidSms, err)
		}

		return err
	}

	newTxn := repository.NewTXNRecord{
		ChatID:      txnMsg.ChatID,
		UserName:    txnMsg.UserName,
		BankName:    msg.BankName,
		CardNumber:  msg.CardNumber,
		TxnDate:     msg.Transaction.Date,
		Type:        msg.Transaction.Type,
		Currency:    msg.Currency,
		Amount:      msg.Amount,
		Payee:       msg.Payee,
		CountryCode: msg.CountryCode,
		SmsMessage:  msg.OriginalMsg,
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}
	err = t.txnRepo.Save(ctx, &newTxn)
	if err != nil {
		return errors.Wrap(err, "can't save transaction to DB")
	}

	return nil
}
