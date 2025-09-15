package config

type Configuration struct {
	BindAddress    string `toml:"bind_address"`
	LogLevel       string `toml:"log_level"`
	FrontOriginUrl string `toml:"front_origin_url"`
	LocalOriginUrl string `toml:"local_origin_url"`
	// tdb
	TarantooldbAddress string `toml:"tarantooldb_address"`
	Username           string `toml:"username"`
	Password           string `toml:"password"`
	Timeout            uint   `toml:"timeout"`
	// kafka
	KafkaAddress string `toml:"kafka_address"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		BindAddress:        ":8080",
		LogLevel:           "debug",
		FrontOriginUrl:     "http://localhost:3000",
		LocalOriginUrl:     "http://localhost:3000",
		TarantooldbAddress: ":3301",
		Username:           "username",
		Password:           "password",
		Timeout:            3,
		KafkaAddress:       ":9092",
	}
}
