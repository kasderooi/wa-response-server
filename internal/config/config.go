package config

import (
	"flag"
	env "github.com/caarlos0/env/v10"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	HttpAddr      string        `env:"HTTP_ADDR" required:"false" envDefault:":8080"`
	HttpIOTimeout time.Duration `env:"TIMEOUT" envDefault:"60s" desc:"timeout for HTTP read and write"`

	DbDsn           string        `env:"DB_DSN" envDefault:"file:whatsapp.db?_foreign_keys=on" required:"true"`
	DbDialect       string        `env:"DB_DIALECT" envDefault:"sqlite3"`
	DbLogLevel      string        `env:"DB_LOG_LEVEL" envDefault:"DEBUG"`
	MaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME" envDefault:"5m"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"2"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"5"`
}

// ParseFromEnvVars - populates config from environment
func ParseFromEnvVars() (Config, error) {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	if len(flag.Args()) != 0 {
		logrus.Error("extra command-line arguments passed which are not valid")
	}

	return cfg, nil
}
