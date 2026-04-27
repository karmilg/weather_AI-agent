package main

import (
	"context"
	"log"

	"github.com/karmilg/weather_AI-agent/config"
	"github.com/karmilg/weather_AI-agent/internal/agent"
	"github.com/karmilg/weather_AI-agent/internal/cli"
	"github.com/karmilg/weather_AI-agent/internal/database"
	"github.com/karmilg/weather_AI-agent/internal/llm"
	"github.com/karmilg/weather_AI-agent/internal/tools"
)

func main() {
	cfg := config.Load()

	if cfg.WeatherAPIKey == "" {
		log.Fatal("❌ WEATHER_API_KEY не задан")
	}

	db, err := database.NewDB(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Pool.Close()

	llmClient := llm.NewClient(cfg.OllamaURL, cfg.LLMModel)
	weatherTool := tools.NewWeatherTool(cfg.WeatherAPIKey, db)
	historyTool := tools.NewHistoryTool(db)
	timeTool := &tools.TimeTool{}
	agentInstance := agent.NewAgent(llmClient, []tools.Tool{weatherTool, historyTool, timeTool})

	cliApp := cli.NewCli(agentInstance)
	if err := cliApp.Run(); err != nil {
		log.Fatal(err)
	}
}
