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
	priorPrefix = "Priorbank"
)

// Prior is an implementation of SMS Parser for Priorbank
type Prior struct {
	log *log.Logger
}

// NewPrior creates a new instance of Prior sms parser
func NewPrior(log *log.Logger) *Prior {
	return &Prior{log: log}
}

func (p *Prior) Parse(text string) (smsparser.SmsMessage, error) {
	if !strings.HasPrefix(text, priorPrefix) {
		return smsparser.SmsMessage{}, errors.New("SMS either not from Priorbank or format changed")
	}

	msg := smsparser.SmsMessage{}
	chunks := strings.Split(text, ". ")
	if len(chunks) < 5 {
		return smsparser.SmsMessage{}, errors.New("SMS either not from Priorbank or format changed")
	}

	msg.BankName = chunks[0]

	msg.CardNumber = p.getCardNumber(chunks[1])
	msg.Transaction.Date = p.getTransactionDate(chunks[1])

	msg.Transaction.Type = p.getTransactionType(chunks[2])
	msg.Amount = p.getAmount(chunks[2])
	msg.Currency = p.getCurrency(chunks[2])

	msg.CountryCode = p.getCountryCode(chunks[3])
	msg.Payee = p.getPayee(chunks[3])

	msg.OriginalMsg = text

	return msg, nil
}

func (p *Prior) getCardNumber(text string) string {
	cardChunks := strings.Split(text, " ")
	return cardChunks[1]
}

func (p *Prior) getTransactionDate(text string) string {
	cardChunks := strings.Split(text, " ")
	return fmt.Sprintf("%s %s", cardChunks[2], cardChunks[3])
}

func (p *Prior) getTransactionType(text string) string {
	transAndAmountChunks := strings.Split(text, " ")
	return transAndAmountChunks[0]
}

func (p *Prior) getAmount(text string) float64 {
	transAndAmountChunks := strings.Split(text, " ")
	amount, err := strconv.ParseFloat(transAndAmountChunks[1], 64)
	if err != nil {
		p.log.Printf("cant parse a line for transaction amount\n\t\tInvalid text:%s. err: %v", text, err)
		return 0
	}
	return amount
}

func (p *Prior) getCurrency(text string) string {
	transAndAmountChunks := strings.Split(text, " ")
	return transAndAmountChunks[2]
}

func (p *Prior) getCountryCode(text string) string {
	transAndAmountChunks := strings.Split(text, " ")
	return transAndAmountChunks[0]
}

func (p *Prior) getPayee(text string) string {
	chunks := strings.Split(text, " ")
	return fmt.Sprintf("%s %s", chunks[1], chunks[2])
}
