package sms

// Parser describes interface for parsing text SMS message from bank and return SmsMessage
type Parser interface {
	Parse(text string) (Message, error)
}
