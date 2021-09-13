package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// TXNRecord is a new transaction to be saved
type TXNRecord struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ChatID      int64              `bson:"chatID"`
	UserName    string             `bson:"userName"`
	BankName    string             `bson:"bankName"`
	CardNumber  string             `bson:"cardNumber"`
	TxnDate     time.Time          `bson:"txnDate"`
	Type        string             `bson:"type"`
	Currency    string             `bson:"currency"`
	Amount      float64            `bson:"amount"`
	Payee       string             `bson:"payee"`
	CountryCode string             `bson:"countryCode"`
	SmsMessage  string             `bson:"smsMessage"`
	DateCreated time.Time          `bson:"dateCreated"`
	DateUpdated time.Time          `bson:"dateUpdated"`
}

// NewTXNRecord is a new transaction to be saved
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

//go:generate sh -c "mockery --inpackage --name TXNer --print > /tmp/mock.tmp && mv /tmp/mock.tmp txner_mock.go"
// TXNer a repo for managing NewTXNRecord
type TXNer interface {
	Save(ctx context.Context, newTxn *NewTXNRecord) error
}

// InvalidSmsRecord represents a full SMS message that was not parsed successfully
type InvalidSmsRecord struct {
	ChatID      int64     `bson:"chatID"`
	UserName    string    `bson:"userName"`
	DateCreated time.Time `bson:"dateCreated"`
	DateUpdated time.Time `bson:"dateUpdated"`
	SmsMessage  string    `bson:"smsMessage"`
}

//go:generate sh -c "mockery --inpackage --name InvalidSMSer --print > /tmp/mock.tmp && mv /tmp/mock.tmp invalidSMSer_mock.go"
// InvalidSMSer a repo for manipulating InvalidSmsRecord.
type InvalidSMSer interface {
	Save(ctx context.Context, invalidSMS *InvalidSmsRecord) error
}
