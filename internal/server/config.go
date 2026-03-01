package server

import (
	"encoding/base64"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type Base64Bytes []byte

func (b *Base64Bytes) UnmarshalText(text []byte) error {
	decoded, err := base64.StdEncoding.DecodeString(string(text))
	*b = decoded
	return err
}

type AppConfig struct {
	BaseURL *url.URL `env:"BASE_URL,required"`
}

type CookieConfig struct {
	HashKey  Base64Bytes `env:"HASH_KEY,required"`
	BlockKey Base64Bytes `env:"BLOCK_KEY,required"`
}

type DatabaseConfig struct {
	ConnectionString string `env:"CONNECTION_STRING,required"`
}

type JWTConfig struct {
	Secret Base64Bytes `env:"SECRET,required"`
}

type LoggerConfig struct {
	Level string `env:"LEVEL" envDefault:"error"`
}

type MailConfig struct {
	SMTP   string `env:"SMTP"`
	APIKey string `env:"API_KEY"`
	From   string `env:"FROM,required"`
}

type Config struct {
	App      AppConfig      `envPrefix:"APP_"`
	Cookie   CookieConfig   `envPrefix:"COOKIE_"`
	Database DatabaseConfig `envPrefix:"DATABASE_"`
	JWT      JWTConfig      `envPrefix:"JWT_"`
	Logger   LoggerConfig   `envPrefix:"LOGGER_"`
	Mail     MailConfig     `envPrefix:"MAIL_"`
}

func NewConfig() (Config, error) {
	config, err := env.ParseAs[Config]()

	return config, err
}
