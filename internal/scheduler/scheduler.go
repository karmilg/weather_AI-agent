package scheduler

import (
    "log"
    "time"

    "github.com/karmilg/weather_AI-agent/internal/database"
    "github.com/karmilg/weather_AI-agent/internal/telegram"
)

type Scheduler struct {
    db  *database.DB
    bot *telegram.Bot
}

func NewScheduler(db *database.DB, bot *telegram.Bot) *Scheduler {
    return &Scheduler{
        db:  db,
        bot: bot,
    }
}

func (s *Scheduler) Start() {
    log.Println("⏰ Планировщик уведомлений запущен")

    ticker := time.NewTicker(1 * time.Minute)
    go func() {
        for range ticker.C {
            s.checkAndSendNotifications()
        }
    }()
}

func (s *Scheduler) checkAndSendNotifications() {
    loc, err := time.LoadLocation("Europe/Moscow")
    if err != nil {
        log.Printf("❌ Ошибка загрузки часового пояса: %v", err)
        loc = time.UTC
    }

    now := time.Now().In(loc)
    currentTime := now.Format("15:04") 

    log.Printf("⏰ Текущее время (МСК): %s", currentTime)

    subs, err := s.db.GetSubscriptionsByTime(currentTime)
    if err != nil {
        log.Printf("❌ Ошибка получения подписок: %v", err)
        return
    }

    if len(subs) == 0 {
        return
    }

    log.Printf("📋 Найдено подписок на %s: %d", currentTime, len(subs))

    for _, sub := range subs {
        log.Printf("📤 Отправка уведомления для %s (chat: %d)", sub.City, sub.ChatID)

        if err := s.bot.SendWeatherNotification(sub.ChatID, sub.City); err != nil {
            log.Printf("❌ Ошибка отправки для %s: %v", sub.City, err)
        } else {
            log.Printf("✅ Уведомление отправлено для %s", sub.City)

			if err := s.db.MarkSubscriptionSent(sub.ID); err != nil {
                log.Printf("⚠️ Не удалось отметить отправку: %v", err)
            }
        }
    }
}