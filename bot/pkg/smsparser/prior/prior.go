package prior

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/oneils/ynab-helper/bot/pkg/smsparser"
)

const (
	priorPrefix         = "Priorbank"
	cardNumberIndx      = 1
	transactionTypeIndx = 0
	transactionDateIndx = 2
	transactionTimeIndx = 3
	currencyIndx        = 2
	countryCodeIndx     = 0
	payeeIndx1          = 1
	payeeIndx2          = 2
	amountIndx          = 1
)

// Prior is an implementation of SMS Parser for Priorbank
type Prior struct {
	log *log.Logger
}

// New creates a new instance of Prior sms parser
func New(log *log.Logger) *Prior {
	return &Prior{log: log}
}

// Parse parses the specified text and returns SmsMessage if there are no errors.
// if a field could not be parsed, it will be set by zero value (empty, 0, etc)
func (p *Prior) Parse(text string) (smsparser.SmsMessage, error) {
	if !strings.HasPrefix(text, priorPrefix) {
		return smsparser.SmsMessage{}, errors.New("SMS either not from Priorbank or format changed")
	}

	msg := smsparser.SmsMessage{}
	chunks := strings.Split(text, ". ")

	// a message from Priorbank should be splited into the following 5 chunks:
	// Priorbank
	// Karta 4***3345 10-09-2021 15:40:19
	// Oplata 38.96 BYN
	// BLR SHOP SOSEDI
	// Spravka: 80171199900
	if len(chunks) < 5 {
		return smsparser.SmsMessage{}, errors.New("SMS either not from Priorbank or format changed")
	}

	msg.BankName = chunks[0]

	// line with card number and transaction date is valid
	// Karta 4***3345 10-09-2021 15:40:19
	if lineValid(chunks[1], 4) {
		msg.CardNumber = p.parseValue(chunks[1], cardNumberIndx)
		msg.Transaction.Date = p.parseCompositeValue(chunks[1], transactionDateIndx, transactionTimeIndx)
	}

	// line with transaction type, amount and currency is valid
	// Oplata 38.96 BYN
	if lineValid(chunks[2], 3) {
		msg.Transaction.Type = p.parseValue(chunks[2], transactionTypeIndx)
		msg.Amount = p.getAmount(chunks[2], amountIndx)
		msg.Currency = p.parseValue(chunks[2], currencyIndx)
	}

	// line with country code and payee name is valid
	// BLR SHOP SOSEDI.
	if lineValid(chunks[3], 2) {

		countryCode, payee := p.getCountryCodeAndPayee(chunks[3], countryCodeIndx)

		msg.CountryCode = countryCode
		msg.Payee = payee
	}

	msg.OriginalMsg = text

	return msg, nil
}

func lineValid(text string, minLength int) bool {
	return len(strings.Split(text, " ")) >= minLength
}

func (p *Prior) parseValue(text string, valueIndex int) string {
	chunks := strings.Split(text, " ")
	if valueIndex > len(chunks) {
		return ""
	}
	return chunks[valueIndex]
}

func (p *Prior) parseCompositeValue(text string, indx1, indx2 int) string {
	chunks := strings.Split(text, " ")
	if indx1 > len(chunks)-1 || indx2 > len(chunks)-1 {
		return ""
	}
	return fmt.Sprintf("%s %s", chunks[indx1], chunks[indx2])
}

func (p *Prior) getAmount(text string, amountIndx int) float64 {
	amountStr := p.parseValue(text, amountIndx)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		p.log.Printf("cant parse a line for transaction amount\n\t\tInvalid text:%s. err: %v", text, err)
		return 0
	}
	return amount
}

func (p *Prior) getCountryCodeAndPayee(text string, countryCodeIndx int) (string, string) {
	chunks := strings.Split(text, " ")
	if len(chunks) <= 1 {
		return "", ""
	}

	code := chunks[countryCodeIndx]
	// verify code looks like BLR
	if len(code) == 3 {
		// paye may have different words amount in the name
		payee := strings.Join(chunks[1:], " ")

		return code, payee
	}
	return "", ""
}
