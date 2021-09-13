package transaction

// TxnMessage is a RAW message from a user to be parsed before saving it DB.
type TxnMessage struct {
	ChatID   int64
	UserName string
	Text     string
}
