package config

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort              string `env:"SERVER_PORT"                envDefault:"8080"`
	LogLevel                string `env:"LOG_LEVEL"                  envDefault:"INFO"`
	DBHost                  string `env:"POSTGRES_HOST"              envDefault:"db"`
	DBPort                  string `env:"DB_PORT"                    envDefault:"5432"`
	DBUser                  string `env:"POSTGRES_USER"              envDefault:"user"`
	DBPassword              string `env:"POSTGRES_PASSWORD"          envDefault:"password"`
	DBName                  string `env:"DB_NAME"                    envDefault:"db"`
	SSLMode                 string `env:"SSLMode"                    envDefault:"disable"`
	GooseDriver             string `env:"GOOSE_DRIVER"               envDefault:"postgres"`
	ReadHeaderTimeoutSecond int    `env:"READ_HEADER_TIMEOUT_SECOND" envDefault:"5"`
	JWTSecretKey            string `env:"JWT_SECRET_KEY"             envDefault:"secret"`
	ExpirationTimeHours     int    `env:"EXPIRATION_TIME_HOURS"     envDefault:"24"`
	TelegramBotToken        string `env:"TELEGRAM_BOT_TOKEN"         envDefault:""`
	TwoFAHost               string `env:"TWO_FA_HOST"                envDefault:"bot"`
	TwoFAPort               string `env:"TWO_FA_PORT"                envDefault:"8080"`
	PPROF                   bool   `env:"PPROF"                      envDefault:"false"`
}

func MustLoad() (cfg *Config, err error) {
	err = godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg = &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	return cfg, nil
}
