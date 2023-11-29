package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Database    string `env:"DATABASE" envDefault:"postgresql" validate:"oneof=postgresql mongodb"`
	PostgresURI string `env:"POSTGRES_DB_URI" validate:"uri"`
	MongoURI    string `env:"MONGO_DB_URI" validate:"uri"`
	RedisURI    string `env:"REDIS_URI" validate:"uri"`
}

func NewConfig() Config {
	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	return conf
}
