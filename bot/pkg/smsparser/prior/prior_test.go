package prior

import (
	"testing"

	"github.com/oneils/ynab-helper/bot/pkg/smsparser"
	"github.com/stretchr/testify/assert"
)

func TestParse_verifyErrors(t *testing.T) {
	testTable := []struct {
		name    string
		smsText string
		errText string
	}{
		{
			name:    "should return error when message does not start from Prior",
			smsText: "some text",
			errText: "SMS either not from Priorbank or format changed",
		},
		{
			name:    "should return error when message is empty",
			smsText: "",
			errText: "SMS either not from Priorbank or format changed",
		},
		{
			name:    "should return error when message was split by a wrong amount of chunks",
			smsText: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. BLR SHOP SOSEDI",
			errText: "SMS either not from Priorbank or format changed",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			prior := Prior{}

			result, err := prior.Parse(testCase.smsText)
			assert.Empty(t, result)
			assert.EqualError(t, err, testCase.errText)
		})
	}
}
func TestParse2(t *testing.T) {

	testTable := []struct {
		name           string
		smsText        string
		expectedSmsMsg smsparser.SmsMessage
	}{
		{
			name:    "should parse message to struct",
			smsText: `Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. BLR SHOP SOSEDI.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "BYN",
				Amount:      38.96,
				CountryCode: "BLR",
				Payee:       "SHOP SOSEDI",
				OriginalMsg: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. BLR SHOP SOSEDI.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "Oplata",
					Date: "10-09-2021 15:40:19",
				},
			},
		},
		{
			name:    "should parse message when transaction date is missed",
			smsText: `Priorbank. Karta 4***3345. Oplata 38.96 BYN. BLR SHOP SOSEDI.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "BYN",
				Amount:      38.96,
				CountryCode: "BLR",
				Payee:       "SHOP SOSEDI",
				OriginalMsg: "Priorbank. Karta 4***3345. Oplata 38.96 BYN. BLR SHOP SOSEDI.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "Oplata",
					Date: "",
				},
			},
		},
		{
			name:    "should parse message to struct when card number is missed",
			smsText: `Priorbank. Karta 10-09-2021 15:40:19. Oplata 38.96 BYN. BLR SHOP SOSEDI.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "",
				Currency:    "BYN",
				Amount:      38.96,
				CountryCode: "BLR",
				Payee:       "SHOP SOSEDI",
				OriginalMsg: "Priorbank. Karta 10-09-2021 15:40:19. Oplata 38.96 BYN. BLR SHOP SOSEDI.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "Oplata",
					Date: "10-09-2021 15:40:19",
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			prior := Prior{}

			result, err := prior.Parse(testCase.smsText)

			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedSmsMsg, result)
		})
	}

}
