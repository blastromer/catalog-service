package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds runtime settings loaded from the environment.
type Config struct {
	ServiceName    string
	HTTPAddr       string
	MongoURI       string
	MongoDatabase  string
	CollectionName string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ShutdownGrace  time.Duration
	QueryTimeout   time.Duration
	DefaultLimit   int64
	MaxLimit       int64
}

// Load reads configuration from the environment.
// An optional .env file is loaded when present; missing files are ignored.
func Load(serviceName, defaultAddr, collection string) (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		ServiceName:    serviceName,
		HTTPAddr:       getEnv("HTTP_ADDR", defaultAddr),
		MongoURI:       os.Getenv("MONGO_URI"),
		MongoDatabase:  getEnv("MONGO_DATABASE", "catalog"),
		CollectionName: getEnv("MONGO_COLLECTION", collection),
		ReadTimeout:    getDuration("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:   getDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
		ShutdownGrace:  getDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
		QueryTimeout:   getDuration("MONGO_QUERY_TIMEOUT", 10*time.Second),
		DefaultLimit:   getInt64("DEFAULT_PAGE_LIMIT", 50),
		MaxLimit:       getInt64("MAX_PAGE_LIMIT", 200),
	}

	if cfg.MongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is required")
	}
	if cfg.MongoDatabase == "" {
		return nil, fmt.Errorf("MONGO_DATABASE is required")
	}
	if cfg.CollectionName == "" {
		return nil, fmt.Errorf("collection name is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func getInt64(key string, fallback int64) int64 {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
