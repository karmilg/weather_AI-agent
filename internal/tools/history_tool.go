package tools

import (
	"fmt"
	"strings"

	"github.com/karmilg/weather_AI-agent/internal/database"
)

type HistoryTool struct {
	db *database.DB
}

func NewHistoryTool(db *database.DB) *HistoryTool {
	return &HistoryTool{
		db: db,
	}
}

func (h *HistoryTool) Name() string {
	return "get_weather_history"
}

func (h *HistoryTool) Description() string {
	return "Получить историю погоды. Формат: 'город:последняя', 'город:вчера', 'город:статистика'"
}

func (h *HistoryTool) Execute(input string) (string, error) {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return "Формат: город:команда (последняя, вчера, статистика)", nil
	}

	city := parts[0]
	command := parts[1]

	switch command {
	case "последняя":
		record, err := h.db.GetLastWeather(city)
		if err != nil {
			return fmt.Sprintf("Нет данных для: %s", city), nil
		}

		return fmt.Sprintf("📋 Последний запрос для %s:\n🌡️ %.1f°C\n☁️ %s\n📅 %s",
            record.City, record.Temperature, record.Description, record.RequestedAt.Format("02.01.2006 15:04")), nil
	
	case "вчера":
		record, err := h.db.GetYesterdayWeather(city)
		if err != nil {
			return fmt.Sprintf("Нет данных для: %s", city), nil
		}

		return fmt.Sprintf("📋 Погода в %s вчера:\n🌡️ %.1f°C\n☁️ %s",
            record.City, record.Temperature, record.Description), nil
	case "статистика":
		stats, err := h.db.GetStats(city)
		if err != nil {
			return fmt.Sprintf("Нет статистики для: %s", city), nil
		}

		 return fmt.Sprintf(
            "📊 Статистика для %s:\n📈 Всего запросов: %d\n🌡️ Средняя: %.1f°C\n❄️ Минимум: %.1f°C\n🔥 Максимум: %.1f°C\n💧 Влажность: %.1f%%",
            city, stats["total_requests"], stats["avg_temperature"], stats["min_temperature"], stats["max_temperature"], stats["avg_humidity"]), nil
	default:
        return "Неизвестно. Доступно: последняя, вчера, статистика", nil

	}
}