package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
)

type envFile string

type Config struct {
	MariaDB MariaDB `mapstructure:"MariaDB"`
	Redis   Redis   `mapstructure:"Redis"`
	JWT     JWT     `mapstructure:"Jwt"`
	Oauth   Oauth   `mapstructure:"Oauth"`
	Email   Email   `mapstructure:"Email"`
	APIKey  APIKey  `mapstructure:"ApiKey"`
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

type JWT struct {
	SignKey             string `env:"JWT_SignKey" envDefault:"1sfkfWjfOkQ8hFhka8"`
	AccessTokenExpired  int    `env:"JWT_A_EXPIRED" envDefault:"15"`    // min
	RefreshTokenExpired int    `env:"JWT_R_EXPIRED" envDefault:"20160"` // min
}

type Oauth struct {
	NaverRedirect      string `env:"NAVER_LOGIN_REDIRECT_URL"`
	NaverClientId      string `env:"NAVER_LOGIN_CLIENT_ID"`
	NaverClientSecret  string `env:"NAVER_LOGIN_CLIENT_SECRET"`
	KaKaoRedirect      string `env:"KAKAO_LOGIN_REDIRECT_URL"`
	KaKaoClientId      string `env:"KAKAO_LOGIN_CLIENT_ID"`
	GoogleRedirect     string `env:"GOOGLE_LOGIN_REDIRECT_URL"`
	GoogleClientId     string `env:"GOOGLE_LOGIN_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_LOGIN_CLIENT_SECRET"`
}

type Email struct {
	Host     string `env:"EMAIL_HOST"`
	Password string `env:"EMAIL_PASSWORD"`
	UserName string `env:"EMAIL_USERNAME"`
}

type APIKey struct {
	Businessman string `env:"API_BUSINESS_MAN"`
}

const (
	DEV  envFile = ".env.dev" // default
	PROD envFile = ".env.prod"
	TEST envFile = ".env.test"
)

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(filePath string) error {
	fileName := envMode()
	url := fmt.Sprintf("%s/%s", filePath, fileName)
	envMap, err := godotenv.Read(url)
	if err != nil {
		return err
	}

	op := env.Options{Environment: envMap}

	data := make(map[string]interface{})
	data["Redis"] = &Redis{}
	data["MariaDB"] = &MariaDB{}
	data["Jwt"] = &JWT{}
	data["Oauth"] = &Oauth{}
	data["Email"] = &Email{}
	data["ApiKey"] = &APIKey{}

	for _, v := range data {
		if err := env.Parse(v, op); err != nil {
			fmt.Println(err)
			return err
		}
	}

	if err := mapstructure.Decode(data, &c); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func envMode() envFile {
	mode := os.Getenv("GO_ENV")
	if mode == "TEST" {
		return TEST
	} else if mode == "PROD" {
		return PROD
	}
	return DEV
}
