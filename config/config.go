package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WeatherAPIKey string
	OllamaURL     string
	LLMModel      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	TelegramToken string
	DatabaseURL string
}

func Load() *Config {
	_ = godotenv.Load() 

	databaseURL := os.Getenv("DATABASE_URL")

	cfg := &Config{
		WeatherAPIKey: getEnv("WEATHER_API_KEY"),
		OllamaURL: getEnv("OLLAMA_URL"),
		LLMModel: getEnv("LLM_MODEL"),
		DatabaseURL: databaseURL,
		TelegramToken: getEnv("TELEGRAM_TOKEN"),
	}

	if databaseURL == "" {
		cfg.DBHost = getEnv("DB_HOST")
		cfg.DBPort = getEnv("DB_PORT")
		cfg.DBUser = getEnv("DB_USER")
		cfg.DBPassword = getEnv("DB_PASSWORD")
		cfg.DBName = getEnv("DB_NAME")
	}

	return cfg
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Ошибка: переменная %s не задана! Укажите её в .env файле", key)
	}

	return val
}
