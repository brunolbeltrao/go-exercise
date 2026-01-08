package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        string
	CacheTTL    time.Duration
	HTTPTimeout time.Duration
}

func FromEnv() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		CacheTTL:    parseDuration(getEnv("CACHE_TTL", "60s"), 60*time.Second),
		HTTPTimeout: parseDuration(getEnv("HTTP_TIMEOUT", "3s"), 3*time.Second),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func parseDuration(s string, def time.Duration) time.Duration {
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	if i, err := strconv.Atoi(s); err == nil && i > 0 {
		return time.Duration(i) * time.Second
	}
	return def
}
