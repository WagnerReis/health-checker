package config

import (
	"log"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Port                   string
	NumMaxWorkers          int
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiration  int
	RefreshTokenExpiration int
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

	accessTokenExpiration, err := strconv.Atoi(getEnvOrThrow("ACCESS_TOKEN_EXPIRATION", "15", true))
	if err != nil {
		log.Fatalf("Invalid ACCESS_TOKEN_EXPIRATION: %v", err)
	}
	refreshTokenExpiration, err := strconv.Atoi(getEnvOrThrow("REFRESH_TOKEN_EXPIRATION", "720", true))
	if err != nil {
		log.Fatalf("Invalid REFRESH_TOKEN_EXPIRATION: %v", err)
	}

	cfgOnce.Do(func() {
		cfg = &Config{
			Port:                   getEnvOrThrow("PORT", "8080", true),
			AccessTokenSecret:      getEnvOrThrow("ACCESS_TOKEN_SECRET", "", true),
			RefreshTokenSecret:     getEnvOrThrow("REFRESH_TOKEN_SECRET", "", true),
			AccessTokenExpiration:  accessTokenExpiration,
			RefreshTokenExpiration: refreshTokenExpiration,
		}
	})
	return cfg
}
