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
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("файл .env не найден")
	}

	return &Config{
		WeatherAPIKey: getEnv("WEATHER_API_KEY"),
		OllamaURL: getEnv("OLLAMA_URL"),
		LLMModel: getEnv("LLM_MODEL"),
		DBHost: getEnv("DB_HOST"),
		DBPort: getEnv("DB_PORT"),
		DBUser: getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName: getEnv("DB_NAME"),
		TelegramToken: getEnv("TELEGRAM_TOKEN"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Ошибка: переменная %s не задана! Укажите её в .env файле", key)
	}

	return val
}
