package main

import (
	"flag"
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

	db := store.New(configuration)
	producer := kafka.NewProducer(configuration)
	server := tapmenu.New(configuration, db, producer)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
