package database

import (
	"context"
	"github.com/oneils/ynab-helper/bot/pkg/config"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func InitDB(cfg *config.Config) (*mongo.Database, error) {
	mongoClient, err := newMongoClient(cfg.Mongo)
	if err != nil {
		return nil, err
	}

	return mongoClient.Database(cfg.Mongo.Name), nil
}

func newMongoClient(cfg config.MongoConfig) (*mongo.Client, error) {

	credential := options.Credential{
		Username: cfg.User,
		Password: cfg.Password,
	}
	clientOpts := options.Client().ApplyURI(cfg.URI).SetAuth(credential)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		return nil, errors.Wrapf(err, "error whicle connecting to MongoDB.\n\tURI:\t%s \n\tUser: %s", cfg.URI, cfg.User)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.Wrapf(err, " can't ping MongoDB.\n\tURI:\t%s \n\tUser: %s", cfg.URI, cfg.User)
	}

	return client, nil
}
