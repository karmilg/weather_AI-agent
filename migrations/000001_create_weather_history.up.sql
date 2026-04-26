CREATE TABLE IF NOT EXISTS weather_history (
    id SERIAL PRIMARY KEY,
    city VARCHAR(255) NOT NULL,
    temperature REAL ,
    feels_like REAL,
    humidity INTEGER ,
    wind_speed REAL ,
    pressure INTEGER,
    description TEXT,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_city ON weather_history (city);
CREATE INDEX IF NOT EXISTS idx_requested_at ON weather_history (requested_at);
CREATE INDEX IF NOT EXISTS idx_city_date ON weather_history (city, requested_at);