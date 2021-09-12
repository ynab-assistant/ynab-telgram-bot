package repository

import (
	"context"
	"time"
)

// NewTXN is a new transaction to be saved
type NewTXNRecord struct {
	ChatID      int64     `bson:"chatID"`
	UserName    string    `bson:"userName"`
	BankName    string    `bson:"bankName"`
	CardNumber  string    `bson:"cardNumber"`
	TxnDate     time.Time `bson:"txnDate"`
	Type        string    `bson:"type"`
	Currency    string    `bson:"currency"`
	Amount      float64   `bson:"amount"`
	Payee       string    `bson:"payee"`
	CountryCode string    `bson:"countryCode"`
	SmsMessage  string    `bson:"smsMessage"`
	DateCreated time.Time `bson:"dateCreated"`
	DateUpdated time.Time `bson:"dateUpdated"`
}

// TxnRepository a repo for managing NewTXNRecord
type TxnRepository interface {
	Save(ctx context.Context, newTxn NewTXNRecord) error
}

// InvalidSmsRecord represents a full SMS message that was not parsed successfully
type InvalidSmsRecord struct {
	ChatID      int64     `bson:"chatID"`
	UserName    string    `bson:"userName"`
	DateCreated time.Time `bson:"dateCreated"`
	DateUpdated time.Time `bson:"dateUpdated"`
	SmsMessage  string    `bson:"smsMessage"`
}

// InvalidSmsRepository a repo for manipulating InvalidSmsRecord.
type InvalidSmsRepository interface {
	Save(ctx context.Context, invalidSMS InvalidSmsRecord) error
}
