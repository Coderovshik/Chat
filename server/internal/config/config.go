package config

import (
	"log"
	"net"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Host        string        `env:"SERVER_HOST" env-default:"localhost"`
	Port        string        `env:"SERVER_PORT" env-default:"8080"`
	Timeout     time.Duration `env:"TIMEOUT" env-required:"true"`
	SigningKey  string        `env:"SIGNING_KEY" env-required:"true"`
	TokenTTL    time.Duration `env:"TOKEN_TTL" env-required:"true"`
	DatabaseURI string        `env:"DATABASE_URI" env-required:"true"`
}

func (c *Config) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func New() *Config {
	const path = "config/.env"

	err := godotenv.Load(path)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
