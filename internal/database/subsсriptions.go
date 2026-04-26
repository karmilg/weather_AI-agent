package database

import (
	"context"
	"fmt"
	"log"
)

func (db *DB) AddSubscription(chatID int64, city, notifyTime string) error {
	if len(notifyTime) == 5 {
		notifyTime = notifyTime + ":00"
	}

	query := `INSERT INTO subscriptions (chat_id, city, notify_time) VALUES ($1, $2, $3::time)
	ON CONFLICT (chat_id, city) DO UPDATE SET
	is_active = true, notify_time = EXCLUDED.notify_time, last_sent = NULL`

	_, err := db.Pool.Exec(db.Ctx, query, chatID, city, notifyTime)
	return err
}

func (db *DB) MarkSubscriptionSent(id int) error {
	query := `UPDATE subscriptions SET last_sent = CURRENT_DATE WHERE id = $1`
	_, err := db.Pool.Exec(db.Ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления last_sent: %w", err)
	}
	return nil
}

func (db *DB) RemoveSubscription(chatID int64, city string) error {
	query := `
	UPDATE subscriptions SET is_active = false WHERE chat_id = $1 AND city = $2
	`

	_, err := db.Pool.Exec(context.Background(), query, chatID, city)
	if err != nil {
		log.Printf("❌ Ошибка удаления подписки: %v", err)
        return err
	}

	return nil
}

func (db *DB) GetSubscriptionsByTime(notifyTime string) ([]Subscription, error) {
	if len(notifyTime) == 5 {
        notifyTime = notifyTime + ":00"
	}
	query := `
	SELECT id, chat_id, city, notify_time, is_active, last_sent, created_at
        FROM subscriptions
        WHERE notify_time = $1::time
          AND is_active = true
          AND (last_sent IS NULL OR last_sent < CURRENT_DATE)
	`

	rows, err := db.Pool.Query(db.Ctx, query, notifyTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Subscription
	for rows.Next() {
		var s Subscription
		err := rows.Scan(&s.ID, &s.ChatID, &s.City, &s.NotifyTime, &s.IsActive, &s.LastSent, &s.CreatedAt)
		if err != nil {
			return nil, err
		}

		subs = append(subs, s)
	}

	return subs, nil
}

func (db *DB) GetUserSubscriptions(chatID int64) ([]Subscription, error) {
	log.Printf("🔍 Запрос подписок для chat_id=%d", chatID)
	query := `
	SELECT id, chat_id, city, notify_time, is_active, created_at FROM subscriptions WHERE chat_id = $1 AND is_active = true
	`

	rows, err := db.Pool.Query(context.Background(), query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Subscription
	for rows.Next() {
		var s Subscription
		err := rows.Scan(&s.ID, &s.ChatID, &s.City, &s.NotifyTime, &s.IsActive, &s.CreatedAt)
		if err != nil {
            log.Printf("❌ Ошибка сканирования: %v", err)
            return nil, err
        }
        subs = append(subs, s)
        log.Printf("✅ Найдена подписка: %+v", s) 
    }
    
    log.Printf("📋 Всего подписок для %d: %d", chatID, len(subs)) 
    return subs, nil
}
