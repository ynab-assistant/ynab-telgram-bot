package repository

import (
	"context"
	"github.com/oneils/ynab-helper/bot/platform/database/databasetest"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestMongoDB(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()

	now := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)

	t.Run("should save new TXN Record", func(t *testing.T) {
		newTxnRecord := NewTXNRecord{
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

		txnRepo := NewMongoTxnRepo(db)
		ctx := context.Background()

		result := txnRepo.Save(ctx, &newTxnRecord)

		// verify docs saved correctly
		assert.NoError(t, result)

		cursor, err := txnRepo.collection.Find(ctx, bson.M{})
		assert.NoError(t, err)
		defer cursor.Close(ctx)

		var txnRecords []TXNRecord
		err = cursor.All(ctx, &txnRecords)
		assert.NoError(t, err)

		assert.True(t, len(txnRecords) == 1)
		savedTxn := txnRecords[0]

		assert.Equal(t, newTxnRecord.ChatID, savedTxn.ChatID)
		assert.Equal(t, newTxnRecord.UserName, savedTxn.UserName)
		assert.Equal(t, newTxnRecord.BankName, savedTxn.BankName)
		assert.Equal(t, newTxnRecord.CardNumber, savedTxn.CardNumber)
		assert.Equal(t, newTxnRecord.TxnDate, savedTxn.TxnDate)
		assert.Equal(t, newTxnRecord.Type, savedTxn.Type)
		assert.Equal(t, newTxnRecord.Currency, savedTxn.Currency)
		assert.Equal(t, newTxnRecord.Amount, savedTxn.Amount)
		assert.Equal(t, newTxnRecord.Payee, savedTxn.Payee)
		assert.Equal(t, newTxnRecord.CountryCode, savedTxn.CountryCode)
		assert.Equal(t, newTxnRecord.SmsMessage, savedTxn.SmsMessage)
		assert.Equal(t, newTxnRecord.DateCreated, savedTxn.DateCreated)
		assert.Equal(t, newTxnRecord.DateUpdated, savedTxn.DateUpdated)
	})

	t.Run("should save invalid SMS", func(t *testing.T) {
		newInvalidSMS := InvalidSmsRecord{
			ChatID:      123,
			UserName:    "userName",
			SmsMessage:  "Some test message to parse",
			DateCreated: now,
			DateUpdated: now,
		}

		invalidSmsRepo := NewMongoInvalidSmsRepo(db)
		ctx := context.Background()

		result := invalidSmsRepo.Save(ctx, &newInvalidSMS)

		// verify docs saved correctly
		assert.NoError(t, result)

		cursor, err := invalidSmsRepo.collection.Find(ctx, bson.M{})
		assert.NoError(t, err)
		defer cursor.Close(ctx)

		var invalidRecords []InvalidSmsRecord
		err = cursor.All(ctx, &invalidRecords)
		assert.NoError(t, err)

		assert.True(t, len(invalidRecords) == 1)
		savedSMS := invalidRecords[0]

		assert.Equal(t, newInvalidSMS.ChatID, savedSMS.ChatID)
		assert.Equal(t, newInvalidSMS.UserName, savedSMS.UserName)
		assert.Equal(t, newInvalidSMS.SmsMessage, savedSMS.SmsMessage)
		assert.Equal(t, newInvalidSMS.DateCreated, savedSMS.DateCreated)
		assert.Equal(t, newInvalidSMS.DateUpdated, savedSMS.DateUpdated)
	})
}
