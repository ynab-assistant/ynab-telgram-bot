package smsparser

// SmsMessage represents SMS message from a bank.
type SmsMessage struct {
	BankName    string
	CardNumber  string
	Transaction Transaction
	Currency    string
	Amount      float64
	Payee       string
	CountryCode string
	OriginalMsg string
}

// Transaction represent transaction
type Transaction struct {
	Date string
	Type string
}

// SMSParser describes interface for parsing text SMS message from bank and return SmsMessage
type SMSParser interface {
	Parse(text string) (SmsMessage, error)
}
