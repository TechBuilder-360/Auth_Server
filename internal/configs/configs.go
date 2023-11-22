package configs

import (
	"go.deanishe.net/env"
	"strings"
)

const (
	PRODUCTION  string = "PRODUCTION"
	DEVELOPMENT string = "DEVELOPMENT"
	SANDBOX     string = "SANDBOX"
)

var Instance *Config

type Config struct {
	AppName       string `env:"APP_NAME"`
	Namespace     string `env:"NAMESPACE"`
	BASEURL       string `env:"BASE_URL"`
	Port          string `env:"PORT"`
	Environment   string `env:"ENVIRONMENT"`
	Secret        string `env:"SECRET"`
	TOKENLIFESPAN uint   `env:"TOKEN_LIFE_SPAN"`

	DbName string `env:"DB_NAME"`
	DbUser string `env:"DB_USER"`
	DbPass string `env:"DB_PASS"`
	DbHost string `env:"DB_HOST"`
	DbPort uint   `env:"DB_PORT"`

	RedisURL          string `env:"REDIS_URL"`
	RedisDB           int    `env:"REDIS_DB"`
	RedisPassword     string `env:"REDIS_PASSWORD"`
	RedisCacheRefresh string `env:"REDIS_CACHE_REFRESH"`

	SendGridAPIKey    *string `env:"SENDGRID_API_KEY"`
	SendGridFromEmail *string `env:"SEND_GRID_FROM_EMAIL"`

	PaperTailAppName *string `env:"PAPER_TAIL_APP_NAME"`
	PaperTailPort    *string `env:"PAPER_TAIL_PORT"`

	CloudinaryName   *string `env:"ENV_CLOUD_NAME"`
	CloudinaryAPIKey *string `env:"ENV_CLOUD_API_KEY"`
	CloudinarySecret *string `env:"ENV_CLOUD_API_SECRET"`
}

func Load() {
	c := &Config{}
	if err := env.Bind(c); err != nil {
		panic(err.Error())
	}
	Instance = c
	return
}

func (c *Config) GetEnv() string {
	return strings.ToUpper(Instance.Environment)
}

func IsSandBox() bool {
	return strings.ToUpper(Instance.Environment) == SANDBOX
}

func IsProduction() bool {
	return strings.ToUpper(Instance.Environment) == PRODUCTION
}
