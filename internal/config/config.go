package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Database     string `env:"DATABASE" envDefault:"postgresql" validate:"oneof=postgresql mongodb"`
	PostgresURI  string `env:"POSTGRES_DB_URI" validate:"uri"`
	MongoURI     string `env:"MONGO_DB_URI" validate:"uri"`
	RedisUserURI string `env:"REDIS_USER_URI" validate:"uri"`
	RedisRftURI  string `env:"REDIS_RFT_URI" validate:"uri"`
}

func NewConfig() Config {
	if err := os.Setenv("DATABASE", "postgresql"); err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	if err := os.Setenv("POSTGRES_DB_URI", "postgres://echopguser:pgpw4echo@localhost:5432/echodb?sslmode=disable"); err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	if err := os.Setenv("MONGO_DB_URI", "mongodb://localhost:6543"); err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	if err := os.Setenv("REDIS_USER_URI", "redis://localhost:6379?protocol=3"); err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	if err := os.Setenv("REDIS_RFT_URI", "redis://localhost:6380?protocol=3"); err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	return conf
}
