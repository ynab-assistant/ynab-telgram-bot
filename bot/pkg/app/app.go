package app

import (
	"context"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oneils/ynab-helper/bot/pkg/config"
	"github.com/oneils/ynab-helper/bot/pkg/sms/prior"
	"github.com/oneils/ynab-helper/bot/pkg/telegram"
	"github.com/oneils/ynab-helper/bot/pkg/transaction"
	"github.com/oneils/ynab-helper/bot/pkg/transaction/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Run starts the application
func Run(configPath string) {
	log := log.New(os.Stdout, "BOT : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal("cant init configuration for the app")
	}

	tgbotapi, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("error while creating Telegram Bot API: ", err)
	}
	tgbotapi.Debug = cfg.Telegram.Debug

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("cant create MongoDB client", err)
	}

	parser := prior.New(log)
	txnRepo := repository.NewTxnRepo(db)
	invalidSmsRepo := repository.NewInvalidSmsRepo(db)
	txn := transaction.New(log, parser, txnRepo, invalidSmsRepo)

	bot := telegram.NewBot(tgbotapi, log, txn)

	if err := bot.Start(); err != nil {
		log.Fatalf(" error while starting the bot: %v", err)
	}
}

func initDB(cfg *config.Config) (*mongo.Database, error) {
	mongoClient, err := NewMongoClient(cfg.Mongo)
	if err != nil {
		return nil, err
	}

	return mongoClient.Database(cfg.Mongo.Name), nil
}

func NewMongoClient(cfg config.MongoConfig) (*mongo.Client, error) {

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
