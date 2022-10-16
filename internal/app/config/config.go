package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var Info = new(Config)

type EnvMode string

const (
	DEV  EnvMode = "../../../configs/.env.dev" // default
	PROD EnvMode = "../../../configs/.env.prod"
	TEST EnvMode = "../../../configs/.env.test"
)

func Load() error {
	modeUrl := envMode()
	envMap, err := godotenv.Read(string(modeUrl))
	if err != nil {
		return err
	}

	op := env.Options{Environment: envMap}
	mariaCfg := new(MariaDB)

	if err := env.Parse(mariaCfg, op); err != nil {
		return err
	}

	Info.MariaDB = *mariaCfg

	return nil
}

func envMode() EnvMode {
	mode := os.Getenv("GO_ENV")

	if mode == "TEST" {
		return TEST
	} else if mode == "PROD" {
		return PROD
	}
	return DEV
}

type Config struct {
	MariaDB MariaDB
}

type MariaDB struct {
	Host     string `env:"MARIADB_HOST"`
	User     string `env:"MARIADB_USER"`
	Port     int    `env:"MARIADB_PORT" envDefault:"3306"`
	Database string `env:"MARIADB_DATABASE"`
	Password string `env:"MARIADB_PASSWORD"`
}
