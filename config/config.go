package config

import (
	"log"
	"os"
	"sync"
)

type Config struct {
	Port          string
	NumMaxWorkers int
}

func getEnvOrThrow(key string, defaultValue string, required bool) string {
	value := os.Getenv(key)
	if value == "" && required {
		if defaultValue != "" {
			return defaultValue
		}
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func LoadConfig() *Config {
	var cfg *Config
	cfgOnce := sync.Once{}

	cfgOnce.Do(func() {
		cfg = &Config{
			Port: getEnvOrThrow("PORT", "8080", true),
		}
	})
	return cfg
}
