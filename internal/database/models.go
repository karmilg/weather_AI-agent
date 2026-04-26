package database

import "time"

type WeatherRecord struct {
	ID          int
	City        string
	Temperature float64
	FeelsLike   float64
	Humidity    int
	WindSpeed   float64
	Pressure    int
	Description string
	RequestedAt time.Time
}

type Subscription struct {
	ID         int
	ChatID     int64
	City       string
	NotifyTime string
	IsActive   bool
	LastSent   *time.Time
	CreatedAt  time.Time
}
