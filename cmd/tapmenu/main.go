package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/alex-pvl/go-tapmenu/internal/app/config"
	"github.com/alex-pvl/go-tapmenu/internal/app/store"
	"github.com/alex-pvl/go-tapmenu/internal/app/tapmenu"
	"github.com/alex-pvl/go-tapmenu/internal/app/tapmenu/kafka"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/tapmenu.toml", "path to config file")
}

func main() {
	flag.Parse()

	configuration := config.NewConfiguration()
	if _, err := toml.DecodeFile(configPath, configuration); err != nil {
		log.Fatal(err)
	}

	logger, err := configureLogger(configuration)
	if err != nil {
		log.Fatal(err)
	}

	db := store.New(configuration, logger)
	logger.Infof("connected to tarantool %s:***@%s", configuration.Username, configuration.TarantoolAddress)
	producer := kafka.NewProducer(configuration)
	logger.Infof("created Kafka producer on %s", configuration.KafkaAddress)
	server := tapmenu.New(configuration, db, producer, logger)
	if err := server.Start(); err != nil {
		logger.Error(err)
	}
}

func configureLogger(configuration *config.Configuration) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(configuration.LogLevel)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(level)
	return logger, nil
}
