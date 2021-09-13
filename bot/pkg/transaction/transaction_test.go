package transaction

import (
	"context"
	"github.com/oneils/ynab-helper/bot/pkg/sms"
	"github.com/oneils/ynab-helper/bot/pkg/transaction/repository"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"os"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	logger := log.New(os.Stdout, "prefix ", 0)

	now := time.Now()
	ctx := context.TODO()

	t.Run("should save transaction", func(t *testing.T) {
		parser := sms.MockParser{}
		txnRepo := repository.MockTXNer{}
		invalidSmsRepo := repository.MockInvalidSMSer{}

		trans := New(logger, &parser, &txnRepo, &invalidSmsRepo)

		smsMsg := sms.Message{
			BankName:   "Test Bank Name",
			CardNumber: "4***112",
			Transaction: struct {
				Date time.Time
				Type string
			}{
				Date: now,
				Type: "Oplata",
			},
			Currency:    "BYN",
			Amount:      10,
			Payee:       "Test Payee",
			CountryCode: "BLR",
			OriginalMsg: "Some test message to parse",
		}
		parser.On("Parse", "Some test message to parse").Return(smsMsg, nil)

		txnRecord := repository.NewTXNRecord{
			ChatID:      123,
			UserName:    "userName",
			BankName:    "Test Bank Name",
			CardNumber:  "4***112",
			TxnDate:     now,
			Type:        "Oplata",
			Currency:    "BYN",
			Amount:      10,
			Payee:       "Test Payee",
			CountryCode: "BLR",
			SmsMessage:  "Some test message to parse",
			DateCreated: now,
			DateUpdated: now,
		}
		txnRepo.On("Save", ctx, &txnRecord).Return(nil)

		txnMsg := TxnMessage{
			ChatID:   123,
			UserName: "userName",
			Text:     "Some test message to parse",
		}

		result := trans.Save(ctx, txnMsg, now)

		assert.NoError(t, result)
	})

	t.Run("should return error when can't save transaction to DB", func(t *testing.T) {
		parser := sms.MockParser{}
		txnRepo := repository.MockTXNer{}
		invalidSmsRepo := repository.MockInvalidSMSer{}

		trans := New(logger, &parser, &txnRepo, &invalidSmsRepo)

		parser.On("Parse", mock.AnythingOfType("string")).Return(sms.Message{}, nil)
		txnRepo.On("Save", ctx, mock.Anything).Return(errors.New("DB error"))

		txnMsg := TxnMessage{
			ChatID:   4567,
			UserName: "test user",
			Text:     "Some test message to parse",
		}

		result := trans.Save(ctx, txnMsg, now)

		assert.Error(t, result)
		assert.EqualError(t, result, "can't save transaction to DB")
	})

	t.Run("should save invalid message to appropriate db", func(t *testing.T) {
		parser := sms.MockParser{}
		txnRepo := repository.MockTXNer{}
		invalidSmsRepo := repository.MockInvalidSMSer{}

		trans := New(logger, &parser, &txnRepo, &invalidSmsRepo)

		invalidRecod := repository.InvalidSmsRecord{
			ChatID:      4567,
			UserName:    "test user",
			DateCreated: now,
			DateUpdated: now,
			SmsMessage:  "Some test message to parse",
		}

		parser.On("Parse", "Some test message to parse").Return(sms.Message{}, errors.New("parsing error"))
		invalidSmsRepo.On("Save", ctx, &invalidRecod).Return(nil)

		txnMsg := TxnMessage{
			ChatID:   4567,
			UserName: "test user",
			Text:     "Some test message to parse",
		}

		result := trans.Save(ctx, txnMsg, now)

		assert.Error(t, result)
		assert.EqualError(t, result, "parsing error")
	})

	t.Run("should return error when invalid message can not be saved to db", func(t *testing.T) {
		parser := sms.MockParser{}
		txnRepo := repository.MockTXNer{}
		invalidSmsRepo := repository.MockInvalidSMSer{}

		trans := New(logger, &parser, &txnRepo, &invalidSmsRepo)

		txnMsg := TxnMessage{
			ChatID:   4567,
			UserName: "test user",
			Text:     "Some test message to parse",
		}

		parser.On("Parse", mock.AnythingOfType("string")).Return(sms.Message{}, errors.New("parsing error"))
		invalidSmsRepo.On("Save", ctx, mock.Anything).Return(errors.New("some error while saving invalid sms"))

		result := trans.Save(ctx, txnMsg, now)

		assert.Error(t, result)
		assert.EqualError(t, result, "can't save invalid SMS to DB")
	})
}
