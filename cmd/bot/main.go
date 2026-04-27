package main

import (
	"context"
	"log"

	"github.com/karmilg/weather_AI-agent/config"
	"github.com/karmilg/weather_AI-agent/internal/agent"
	"github.com/karmilg/weather_AI-agent/internal/database"
	"github.com/karmilg/weather_AI-agent/internal/llm"
	"github.com/karmilg/weather_AI-agent/internal/scheduler"
	"github.com/karmilg/weather_AI-agent/internal/telegram"
	"github.com/karmilg/weather_AI-agent/internal/tools"
	"github.com/karmilg/weather_AI-agent/internal/weather"
)

func main() {
	cfg := config.Load()
	if cfg.TelegramToken == "" {
		log.Fatal("❌ TELEGRAM_TOKEN не задан")
	}

	if cfg.WeatherAPIKey == "" {
		log.Fatal("❌ WEATHER_API_KEY не задан")
	}

	db, err := database.NewDB(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Pool.Close()

	llmClient := llm.NewClient(cfg.OllamaURL, cfg.LLMModel)
	weatherClient := weather.NewClient(cfg.WeatherAPIKey)
	weatherTool := tools.NewWeatherTool(cfg.WeatherAPIKey, db)
	historyTool := tools.NewHistoryTool(db)
	timeTool := &tools.TimeTool{}

	agentInstance := agent.NewAgent(llmClient, []tools.Tool{weatherTool, historyTool, timeTool})

	bot, err := telegram.NewBot(cfg.TelegramToken, agentInstance, db, weatherClient)
	if err != nil {
		log.Fatal(err)
	}

	scheduler := scheduler.NewScheduler(db, bot)
	scheduler.Start()
	log.Println("🚀 Бот запущен")
    if err := bot.Start(); err != nil {
        log.Fatal(err)
    }
}
