package prior

import (
	"log"
	"os"
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
func TestParse_happyPath(t *testing.T) {

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
				CardNumber:  "",
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
					Date: "",
				},
			},
		},
		{
			name:    "should parse message to struct when transaction type is missed",
			smsText: `Priorbank. Karta 4***3345 10-09-2021 15:40:19. 38.96 BYN. BLR SHOP SOSEDI.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "",
				Amount:      0,
				CountryCode: "BLR",
				Payee:       "SHOP SOSEDI",
				OriginalMsg: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. 38.96 BYN. BLR SHOP SOSEDI.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "",
					Date: "10-09-2021 15:40:19",
				},
			},
		},

		{
			name:    "should parse message to struct when amount is missed",
			smsText: `Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata BYN. BLR SHOP SOSEDI.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "",
				Amount:      0,
				CountryCode: "BLR",
				Payee:       "SHOP SOSEDI",
				OriginalMsg: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata BYN. BLR SHOP SOSEDI.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "",
					Date: "10-09-2021 15:40:19",
				},
			},
		},
		{
			name:    "should parse message to struct when transaction currency is missed",
			smsText: `Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96. BLR SHOP SOSEDI.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "",
				Amount:      0,
				CountryCode: "BLR",
				Payee:       "SHOP SOSEDI",
				OriginalMsg: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96. BLR SHOP SOSEDI.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "",
					Date: "10-09-2021 15:40:19",
				},
			},
		},
		{
			name:    "should parse message to struct when payee contry code is missed",
			smsText: `Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. SHOP SOSEDI.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "BYN",
				Amount:      38.96,
				CountryCode: "",
				Payee:       "",
				OriginalMsg: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. SHOP SOSEDI.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "Oplata",
					Date: "10-09-2021 15:40:19",
				},
			},
		},
		{
			name:    "should parse message to struct when payee contry code and payee has 1 word is missed",
			smsText: `Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. SHOP.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "BYN",
				Amount:      38.96,
				CountryCode: "",
				Payee:       "",
				OriginalMsg: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. SHOP.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "Oplata",
					Date: "10-09-2021 15:40:19",
				},
			},
		},
		{
			name:    "should parse message to struct when payee name is missed",
			smsText: `Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. BLR.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "BYN",
				Amount:      38.96,
				CountryCode: "",
				Payee:       "",
				OriginalMsg: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. BLR.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "Oplata",
					Date: "10-09-2021 15:40:19",
				},
			},
		},
		{
			name:    "should parse message to struct when payee name has one word",
			smsText: `Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. NLD Yandex.Taxi.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "BYN",
				Amount:      38.96,
				CountryCode: "NLD",
				Payee:       "Yandex.Taxi",
				OriginalMsg: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. NLD Yandex.Taxi.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "Oplata",
					Date: "10-09-2021 15:40:19",
				},
			},
		},
		{
			name:    "should parse message to struct when transaction amount is invalid number",
			smsText: `Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata invalid BYN. BLR SHOP SOSEDI.   Spravka: 80171199900`,
			expectedSmsMsg: smsparser.SmsMessage{
				BankName:    "Priorbank",
				CardNumber:  "4***3345",
				Currency:    "BYN",
				Amount:      0,
				CountryCode: "BLR",
				Payee:       "SHOP SOSEDI",
				OriginalMsg: "Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata invalid BYN. BLR SHOP SOSEDI.   Spravka: 80171199900",
				Transaction: smsparser.Transaction{
					Type: "Oplata",
					Date: "10-09-2021 15:40:19",
				},
			},
		},
	}

	log := log.New(os.Stdout, "TEST : ", 0)
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			prior := Prior{log}

			result, err := prior.Parse(testCase.smsText)

			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedSmsMsg, result)
		})
	}

}

func TestParse_factoryFunction(t *testing.T) {
	log := log.New(os.Stdout, "prefix ", 0)
	expectedPrior := &Prior{log}

	p := NewPrior(log)

	assert.Equal(t, expectedPrior, p)
}

func Test_parseValue(t *testing.T) {
	log := log.New(os.Stdout, "prefix ", 0)
	p := NewPrior(log)

	result := p.parseValue("text_without_spaces", 2)

	assert.Empty(t, result)
}

func Test_parseCompositeValue(t *testing.T) {
	log := log.New(os.Stdout, "prefix ", 0)
	p := NewPrior(log)

	result := p.parseCompositeValue("text_without_spaces", 2, 3)

	assert.Empty(t, result)
}
func Test_getCountryCodeAndPayee(t *testing.T) {
	log := log.New(os.Stdout, "prefix ", 0)
	p := NewPrior(log)

	code, payee := p.getCountryCodeAndPayee("text", 4)

	assert.Empty(t, code)
	assert.Empty(t, payee)
}
