package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Database    string `env:"DATABASE" envDefault:"postgresql"`
	PostgresURL string `env:"POSTGRES_DB_URL"`
	MongoURL    string `env:"MONGO_DB_URL"`
}

func NewConfig() Config {
	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	return conf
}
