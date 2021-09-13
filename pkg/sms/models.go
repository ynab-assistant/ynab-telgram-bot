package sms

import "time"

// Message represents SMS message from a bank.
type Message struct {
	BankName    string
	CardNumber  string
	Transaction struct {
		Date time.Time
		Type string
	}
	Currency    string
	Amount      float64
	Payee       string
	CountryCode string
	OriginalMsg string
}
