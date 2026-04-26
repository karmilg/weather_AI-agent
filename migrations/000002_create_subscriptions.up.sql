CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    city VARCHAR(255) NOT NULL,
    notify_time TIME NOT NULL,
    is_active BOOLEAN DEFAULT false,
    last_sent DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_chat_city UNIQUE(chat_id, city)
);

CREATE INDEX IF NOT EXISTS idx_subsriptions_chat ON subscriptions(chat_id);
CREATE INDEX IF NOT EXISTS idx_subsriptions_time ON subscriptions(notify_time);
CREATE INDEX IF NOT EXISTS idx_subsriptions_active ON subscriptions(is_active);