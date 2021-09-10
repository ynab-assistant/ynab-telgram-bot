package prior

import (
	"testing"

	"github.com/oneils/ynab-helper/bot/pkg/smsparser"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("should return error when message does not start from Prior", func(t *testing.T) {
		prior := Prior{}

		result, err := prior.Parse("some text")

		assert.Empty(t, result)
		assert.EqualError(t, err, "SMS either not from Priorbank or format changed")
	})

	t.Run("should return error when message is empty", func(t *testing.T) {
		prior := Prior{}

		result, err := prior.Parse("")

		assert.Empty(t, result)
		assert.EqualError(t, err, "SMS either not from Priorbank or format changed")
	})

	t.Run("should parse message to struct", func(t *testing.T) {
		msg := `Priorbank. Karta 4***3345 10-09-2021 15:40:19. Oplata 38.96 BYN. BLR SHOP SOSEDI.   Spravka: 80171199900`

		expectedResult := smsparser.SmsMessage{
			BankName:    "Priorbank",
			CardNumber:  "4***3345",
			Currency:    "BYN",
			Amount:      38.96,
			CountryCode: "BLR",
			Payee:       "SHOP SOSEDI",
			OriginalMsg: msg,
			Transaction: smsparser.Transaction{
				Type: "Oplata",
				Date: "10-09-2021 15:40:19",
			},
		}

		prior := Prior{}

		result, err := prior.Parse(msg)

		assert.Nil(t, err)
		assert.Equal(t, expectedResult, result)
	})
}
