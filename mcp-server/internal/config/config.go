package config

import (
	"os"
	"strings"
)

type Config struct {
	APIBaseURL    string
	DefaultUserID string
}

func Load() Config {
	apiBaseURL := strings.TrimRight(os.Getenv("SREGEP_API_BASE_URL"), "/")
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:8080"
	}
	return Config{
		APIBaseURL:    apiBaseURL,
		DefaultUserID: os.Getenv("SREGEP_DEFAULT_USER_ID"),
	}
}
