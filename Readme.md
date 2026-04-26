# 🌤️ Погодный AI-агент с PostgreSQL и Telegram

Полноценный AI-агент на Go с интеграцией LLM, погодным API, базой данных и Telegram-уведомлениями.

## 🚀 Как запустить программу?

```bash
# 1. Запустить Docker
docker-compose up -d

# 2. Скачать модель
docker exec -it ollama ollama pull qwen2.5:7ba

# 3. Настроить .env
cp .env
# Добавить WEATHER_API_KEY и TELEGRAM_TOKEN

# 4. Запустить CLI версию
make run-cli

# Или Telegram бота
make run-bot