package weather

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

type Client struct {
    apiKey string
    http   *http.Client
}

func NewClient(apiKey string) *Client {
    return &Client{
        apiKey: apiKey,
        http:   &http.Client{Timeout: 10 * time.Second},
    }
}

func (c *Client) GetWeather(city string) (*WeatherResponse, error) {
    url := fmt.Sprintf(
        "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ru",
        city, c.apiKey,
    )

    log.Printf("🌐 Запрос к OpenWeatherMap: %s", url)

    resp, err := c.http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API вернул статус %d", resp.StatusCode)
    }

    var result WeatherResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    log.Printf("✅ OpenWeatherMap ответил: %s, %.1f°C", result.Name, result.Main.Temp)
    return &result, nil
}