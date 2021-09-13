package sms

//go:generate sh -c "mockery --inpackage --name Parser --print > /tmp/mock.tmp && mv /tmp/mock.tmp parser_mock.go"
// Parser describes interface for parsing text SMS message from bank and return SmsMessage
type Parser interface {
	Parse(text string) (Message, error)
}
