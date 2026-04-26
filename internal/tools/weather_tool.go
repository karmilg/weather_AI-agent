package tools

import (
	"github.com/karmilg/weather_AI-agent/internal/database"
	"github.com/karmilg/weather_AI-agent/internal/weather"
)

type WeatherTool struct {
	weatherClient *weather.Client
	db            *database.DB
}

func NewWeatherTool(apiKey string, db *database.DB) *WeatherTool {
	return &WeatherTool{
		weatherClient: weather.NewClient(apiKey),
		db:            db,
	}
}

func (w *WeatherTool) Name() string {
	return "get_weather"
}

func (w *WeatherTool) Description() string {
	return "Получить текущую погоду для города"
}

func (w *WeatherTool) Execute(input string) (string, error) {
	if input == "" {
		return "Укажите город", nil
	}

	weatherData, err := w.weatherClient.GetWeather(input)
	if err != nil {
		return err.Error(), nil
	}

	temp := weatherData.Main.Temp
	var advice string
	switch {
	case temp <= -10:
        advice = "❄️ Очень холодно! Одевайся максимально тепло."
    case temp <= 0:
        advice = "🧥 Холодно. Надевай шапку и перчатки."
    case temp <= 10:
        advice = "🍂 Прохладно. Куртка или пальто."
    case temp <= 20:
        advice = "🌤️ Тепло. Легкая куртка или кофта."
    default:
        advice = "☀️ Жарко! Футболка и шорты."
    }

	record := database.WeatherRecord{
		City:        weatherData.Name,
        Temperature: weatherData.Main.Temp,
        FeelsLike:   weatherData.Main.FeelsLike,
        Humidity:    weatherData.Main.Humidity,
        WindSpeed:   weatherData.Wind.Speed,
        Pressure:    weatherData.Main.Pressure,
        Description: weatherData.Weather[0].Description,
	}
	w.db.SaveWeather(record)

	return weather.FormatWeather(weatherData) + "\n\n💡 Совет: " + advice, nil
}
