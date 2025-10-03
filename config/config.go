package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	AppEnv     string
	AppVersion string
	BaseURL    string
	DBDsn      string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		AppPort:    getenv("APP_PORT", "4002"),
		AppEnv:     getenv("APP_ENV", "development"),
		AppVersion: getenv("APP_VERSION", "dev"),
		BaseURL:    getenv("BASE_URL", "https://api.gpadaka.com"),
		DBDsn:      must("DATABASE_DSN"),
	}
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(k + " belum di-set")
	}
	return v
}

var AppStartTime = time.Now()
