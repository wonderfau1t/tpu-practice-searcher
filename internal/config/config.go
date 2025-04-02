package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env              string `yaml:"env"`
	AccessSecret     string `env:"JWT_SECRET_TOKEN" env-required:"true"`
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
	Storage          `yaml:"storage"`
	HTTPServer       `yaml:"http_server"`
}

type Storage struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Dbname   string `yaml:"dbname" env-required:"true"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost"`
	Port    int    `yaml:"port" env-default:"8080"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("env variable CONFIG_PATH is not set")
	}

	accessSecret := os.Getenv("JWT_SECRET_TOKEN")
	if accessSecret == "" {
		log.Fatalf("env variable JWT_SECRET_TOKEN is not set")
	}

	refreshSecret := os.Getenv("TELEGRAM_BOT_TOKEN")
	if refreshSecret == "" {
		log.Fatalf("env variable REFRESH_SECRET is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file is not exists: %s", err)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}

	return &config
}
