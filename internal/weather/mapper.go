package weather

import "fmt"

func FormatWeather(w *WeatherResponse) string {
	if w == nil {
		return "Данные о погоде не получены!"
	}

	return fmt.Sprintf(
		 "🌍 Город: %s\n🌡️ Температура: %.1f°C (ощущается %.1f°C)\n💧 Влажность: %d%%\n🌬️ Ветер: %.1f м/с\n📈 Давление: %d мм рт.ст.\n☁️ %s",
		 w.Name, w.Main.Temp, w.Main.FeelsLike, w.Main.Humidity, w.Wind.Speed, w.Main.Pressure, w.Weather[0].Description,
	)
}
