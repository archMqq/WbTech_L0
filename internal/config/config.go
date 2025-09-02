package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServer `yaml:"http_server"`
	Postgre    PostgreCfg `yaml:"postgresql"`
	Kafka      KafkaCfg   `yaml:"kafka"`
}

type HttpServer struct {
	Port string `yaml:"port" env-default:":8080"`
}

type PostgreCfg struct {
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db_name"`
}

type KafkaCfg struct {
	Broker string `yaml:"broker"`
	Topic  string `yaml:"topic"`
}

func Load() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is unsetted")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		log.Fatalf("error opening config: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return &cfg
}
