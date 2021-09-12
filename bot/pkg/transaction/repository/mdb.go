package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	txnCollection        = "transactions"
	invalidSMSCollection = "invalidSMS"
)

// TxnRepo is MongoDB repository for working with Transaction
type TxnRepo struct {
	collection *mongo.Collection
}

// NewTxnRepo creates a new TxnRepo
func NewTxnRepo(db *mongo.Database) *TxnRepo {
	return &TxnRepo{collection: db.Collection(txnCollection)}
}

// Save store the specified TXB to the DB
func (t *TxnRepo) Save(ctx context.Context, newTxn NewTXNRecord) error {
	_, err := t.collection.InsertOne(ctx, newTxn)
	if err != nil {
		return err
	}

	return nil
}

type InvalidSmsRepo struct {
	collection *mongo.Collection
}

func NewInvalidSmsRepo(db *mongo.Database) *InvalidSmsRepo {
	return &InvalidSmsRepo{
		collection: db.Collection(invalidSMSCollection),
	}
}

func (i *InvalidSmsRepo) Save(ctx context.Context, invalidSMS InvalidSmsRecord) error {
	_, err := i.collection.InsertOne(ctx, invalidSMS)
	if err != nil {
		return err
	}

	return nil
}
