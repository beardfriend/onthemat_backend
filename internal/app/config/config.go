package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
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

	data := make(map[string]interface{})
	data["Redis"] = &Redis{}
	data["MariaDB"] = &MariaDB{}

	for _, v := range data {
		if err := env.Parse(v, op); err != nil {
			return err
		}
	}

	if err := mapstructure.Decode(data, &Info); err != nil {
		fmt.Println(err)
	}

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
	MariaDB MariaDB `mapstructure:"MariaDB"`
	Redis   Redis   `mapstructure:"Redis"`
}

type MariaDB struct {
	Host     string `env:"MARIADB_HOST"`
	User     string `env:"MARIADB_USER"`
	Port     int    `env:"MARIADB_PORT" envDefault:"3306"`
	Database string `env:"MARIADB_DATABASE"`
	Password string `env:"MARIADB_PASSWORD"`
}

type Redis struct {
	Host string `env:"REDIS_HOST" envDefault:"localhost"`
	Port int    `env:"REDIS_PORT" envDefault:"6379"`
}
