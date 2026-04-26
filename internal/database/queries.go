package database

import (
	"context"
	"fmt"
)

func (db *DB) SaveWeather(record WeatherRecord) error {
	query := `INSERT INTO weather_history 
	(city, temperature, feels_like, humidity, wind_speed, pressure, description) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Pool.Exec(
		context.Background(),
		query,
		record.City,
		record.Temperature,
		record.FeelsLike,
		record.Humidity,
		record.WindSpeed,
		record.Pressure,
		record.Description,
	)

	return err
}

func (db *DB) GetLastWeather(city string) (*WeatherRecord, error) {
	query := `SELECT city, temperature, feels_like, humidity, wind_speed,
	pressure, description, requested_at FROM weather_history
	WHERE city = $1 
	ORDER_BY requested_at DESC LIMIT 1`

	var record WeatherRecord
	err := db.Pool.QueryRow(db.Ctx, query, city).Scan(
		&record.City, &record.Temperature, &record.FeelsLike, &record.Humidity, &record.WindSpeed,
		&record.Pressure, &record.Description, &record.RequestedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("Нет данных для %s", city)
	}

	return &record, nil
}

func (db *DB) GetYesterdayWeather(city string) (*WeatherRecord, error) {
	query := `SELECT city, temperature, feels_like, humidity, wind_speed,
	pressure, description, requested_at FROM weather_history
	WHERE city = $1 AND requested_at::DATE = CURRENT_DATE - 1
	ORDER_BY requested_at DESC LIMIT 1`

	var record WeatherRecord
	err := db.Pool.QueryRow(db.Ctx, query, city).Scan(
		&record.City, &record.Temperature, &record.FeelsLike, &record.Humidity, &record.WindSpeed,
		&record.Pressure, &record.Description, &record.RequestedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("Нет данных за вчера для %s", city)
	}

	return &record, nil
}

func (db *DB) GetHistory(city string, limit int) ([]WeatherRecord, error) {
	query := `SELECT city, temperature, feels_like, humidity, wind_speed,
	pressure, description, requested_at FROM weather_history
	WHERE city = $1 
	ORDER_BY requested_at DESC LIMIT $2`

	rows, err := db.Pool.Query(db.Ctx, query, city, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []WeatherRecord
	for rows.Next() {
		var r WeatherRecord
		err := rows.Scan(&r.City, &r.Temperature, &r.FeelsLike, &r.Humidity, &r.WindSpeed,
			&r.Pressure, &r.Description, &r.RequestedAt)
		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}

func (db *DB) GetStats(city string) (map[string]interface{}, error) {
	query := `
	SElECT 
	COUNT(*) as total_requests,
	COALESCE(AVG(temperature), 0) as avg_temp,
	COALESCE(MIN(temperature), 0) as min_temp,
	COALESCE(MAX(temperature), 0) as max_temp,
	COALESCE(AVG(humidity), 0) as avg_humidity,
	FROM weather_history
	WHERE city = $1
	`

	var total int
	var avgTemp, minTemp, maxTemp, avgHumidity float64

	err := db.Pool.QueryRow(db.Ctx, query, city).Scan(
		&total, &avgTemp, &minTemp, &maxTemp, &avgHumidity,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_requests": total,
		"avg_temp":       avgTemp,
		"min_temp":       minTemp,
		"max_temp":       maxTemp,
		"avg_humidity":   avgHumidity,
	}, nil
}
