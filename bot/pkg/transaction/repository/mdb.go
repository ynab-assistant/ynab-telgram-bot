package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	txnCollection        = "transactions"
	invalidSMSCollection = "invalidSMS"
)

// MongoTxnRepo is MongoDB repository for working with Transaction
type MongoTxnRepo struct {
	collection *mongo.Collection
}

// NewMongoTxnRepo creates a new TxnRepo
func NewMongoTxnRepo(db *mongo.Database) *MongoTxnRepo {
	return &MongoTxnRepo{collection: db.Collection(txnCollection)}
}

// Save store the specified TXB to the DB
func (t *MongoTxnRepo) Save(ctx context.Context, newTxn NewTXNRecord) error {
	_, err := t.collection.InsertOne(ctx, newTxn)
	if err != nil {
		return err
	}

	return nil
}

// MongoInvalidSmsRepo is a repository for manipulating of invalid SMS messages that could not be parsed
type MongoInvalidSmsRepo struct {
	collection *mongo.Collection
}

// NewMongoInvalidSmsRepo creates a new InvalidSmsRepo
func NewMongoInvalidSmsRepo(db *mongo.Database) *MongoInvalidSmsRepo {
	return &MongoInvalidSmsRepo{collection: db.Collection(invalidSMSCollection)}
}

// Save stores invalid sms that were not parsed or processed correctly
func (i *MongoInvalidSmsRepo) Save(ctx context.Context, invalidSMS *InvalidSmsRecord) error {
	_, err := i.collection.InsertOne(ctx, invalidSMS)
	if err != nil {
		return err
	}

	return nil
}
