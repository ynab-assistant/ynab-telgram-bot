package config

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	defaultMogoDBName    = "ynab-helper"
	defaultMongoURI      = "mongodb://mongo:27017"
	defaultMongoUser     = "root"
	defaultMongoPassword = "root"
)

// Config is application configuration
type Config struct {
	Telegram TelegramConfig
	Mongo    MongoConfig
}

// MongoConfig is a configuration for connecting to a MongoDB
type MongoConfig struct {
	Name     string `mapstructure:"name"`
	URI      string `mapstructure:"uri"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// TelegramConfig is a configuration for Telegram Bot API
type TelegramConfig struct {
	Token string `mapstructure:"token"`
	Debug bool   `mapstructure:"debug"`
}

// Init initialize Config struct
func Init(configPath string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configPath); err != nil {
		return nil, errors.Wrapf(err, "config.Init() ")
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}

func populateDefaults() {
	viper.SetDefault("mongo.name", defaultMogoDBName)
	viper.SetDefault("mongo.uri", defaultMongoURI)
	viper.SetDefault("mongo.user", defaultMongoUser)
	viper.SetDefault("mongo.password", defaultMongoPassword)
}

func parseConfigFile(filePath string) error {

	// TODO: should parse file from any specified palce

	path := strings.Split(filePath, "/")

	// TODO: validation should be added in order do not panic
	viper.AddConfigPath(path[0])
	viper.SetConfigName(path[1])

	return viper.ReadInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("telegram", &cfg.Telegram); err != nil {
		return errors.Wrap(err, "config.unmarshal() telegram ")
	}

	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return errors.Wrap(err, "config.unmarshal() mongoDB ")
	}
	return nil
}
