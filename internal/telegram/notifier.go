package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) SendWeatherNotification(chatID int64, city string) error {
	weatherData, err := b.weather.GetWeather(city)
	if err != nil {
		return err
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

	text := fmt.Sprintf( 
		"🌤️ Добрый день! Прогноз для %s на сегодня:\n\n"+
        "🌡️ %.1f°C (ощущается %.1f°C)\n"+
        "💧 Влажность: %d%%\n"+
        "🌬️ Ветер: %.1f м/с\n"+
		" \n\n💡 Совет: %s"+
        "☁️ %s\n\nХорошего дня! ☀️", 
		weatherData.Name,
        weatherData.Main.Temp,
        weatherData.Main.FeelsLike,
        weatherData.Main.Humidity,
        weatherData.Wind.Speed,
		advice,
        weatherData.Weather[0].Description,
	)

	msg := tgbotapi.NewMessage(chatID, text)
	_, err = b.api.Send(msg)
	return err
}