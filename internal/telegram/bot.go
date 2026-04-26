package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/karmilg/weather_AI-agent/internal/agent"
	"github.com/karmilg/weather_AI-agent/internal/database"
	"github.com/karmilg/weather_AI-agent/internal/weather"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	agent   *agent.Agent
	db      *database.DB
	weather *weather.Client
}

func NewBot(token string, agentInstance *agent.Agent, db *database.DB, weatherClient *weather.Client) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = false
	log.Printf("Бот авторизован как: %s", api.Self.UserName)

	return &Bot{
		api:     api,
		agent:   agentInstance,
		db:      db,
		weather: weatherClient,
	}, nil
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		go b.handleMessage(update.Message)
	}

	return nil
}
