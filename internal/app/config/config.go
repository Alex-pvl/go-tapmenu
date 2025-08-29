package config

type Configuration struct {
	BindAddress string `toml:"bind_address"`
	LogLevel    string `toml:"log_level"`
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
		TarantooldbAddress: ":3301",
		Username:           "username",
		Password:           "password",
		Timeout:            3,
		KafkaAddress:       ":9092",
	}
}
