package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env              string `env:"environment" env-default:"local"`
	AccessSecret     string `env:"JWT_SECRET_TOKEN" env-required:"true"`
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
	Storage
}

type Storage struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Dbname   string `env:"POSTGRES_DB" env-required:"true"`
}

func MustLoad() *Config {
	//configPath := os.Getenv("CONFIG_PATH")
	//if configPath == "" {
	//	log.Fatalf("env variable CONFIG_PATH is not set")
	//}

	//accessSecret := os.Getenv("JWT_SECRET_TOKEN")
	//if accessSecret == "" {
	//	log.Fatalf("env variable JWT_SECRET_TOKEN is not set")
	//}
	//
	//refreshSecret := os.Getenv("TELEGRAM_BOT_TOKEN")
	//if refreshSecret == "" {
	//	log.Fatalf("env variable TELEGRAM_BOT_TOKEN is not set")
	//}
	//
	//if _, err := os.Stat(configPath); err != nil {
	//	log.Fatalf("config file is not exists: %s", err)
	//}

	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("cannot read config from env: %s", err)
	}

	return &config
}
