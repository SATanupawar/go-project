package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  struct {
		Addr string `yaml:"addr"`
	} `yaml:"http_server" env-required:"true"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configFlag := flag.String("config", "", "path to config file")
		flag.Parse()
		configPath = *configFlag

		if configPath == "" {
			// Look for default config locations instead of failing immediately
			defaultLocations := []string{
				"local.yaml",
				"config/local.yaml",
				"configs/local.yaml",
			}

			for _, loc := range defaultLocations {
				if _, err := os.Stat(loc); err == nil {
					configPath = loc
					break
				}
			}

			if configPath == "" {
				log.Fatalf("CONFIG_PATH is not set and no config file found in default locations")
			}
		}
	}

	// Print the path being used to help with debugging
	log.Printf("Using config file: %s", configPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}

	// Validate required fields
	if cfg.StoragePath == "" {
		log.Fatalf("storage_path is required in config")
	}

	// Validate HTTP server address
	if cfg.HttpServer.Addr == "" {
		log.Fatalf("http_server.addr is required in config")
	}

	return &cfg
}
