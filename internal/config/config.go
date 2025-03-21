package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	HttpServer  struct {
		Address string `yaml:"address"`
	} `yaml:"http_server"`
}

func MustLoad() *Config {

	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {

		flags := flag.String("config", "", "path to config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatalf("CONFIG_PATH is not set")
		}

	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}

	return &cfg

}
