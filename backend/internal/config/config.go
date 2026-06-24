package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	SupabaseURL    string
	SupabaseAPIKey string
	Port           string
	AllowedOrigin  string
}

func Load() (Config, error) {
	cfg := Config{
		SupabaseURL:    strings.TrimRight(os.Getenv("SUPABASE_URL"), "/"),
		SupabaseAPIKey: os.Getenv("SUPABASE_API_KEY"),
		Port:           valueOrDefault(os.Getenv("PORT"), "8080"),
		AllowedOrigin:  valueOrDefault(os.Getenv("ALLOWED_ORIGIN"), "http://localhost:3000"),
	}

	if cfg.SupabaseURL == "" {
		return Config{}, errors.New("SUPABASE_URL is required")
	}
	if cfg.SupabaseAPIKey == "" {
		return Config{}, errors.New("SUPABASE_API_KEY is required")
	}

	return cfg, nil
}

func valueOrDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
