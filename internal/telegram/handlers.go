package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := msg.Text

	waitingMsg := tgbotapi.NewMessage(chatID, "⏳ Подождите, получаю информацию...")
	sentMsg, _ := b.api.Send(waitingMsg)

	if strings.HasPrefix(text, "/") {
		b.handleCommand(chatID, text)
		return
	}

	answer, err := b.agent.Run(text)
	if err != nil {
		answer = "Ошибка: " + err.Error()
	}

	deleteMsg := tgbotapi.NewDeleteMessage(chatID, sentMsg.MessageID)
	b.api.Send(deleteMsg)

	reply := tgbotapi.NewMessage(chatID, answer)
	b.api.Send(reply)
}

func (b *Bot) handleCommand(chatID int64, command string) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "/start":
		msg := tgbotapi.NewMessage(chatID,
			"🌤️ Привет! Я погодный AI-агент.\n\n"+
				"📌 Что я умею:\n"+
				"• Спроси погоду: 'какая погода в Москве?'\n"+
				"• /subscribe Москва 08:00 — подписка\n"+
				"• /unsubscribe Москва — отписка\n"+
				"• /mysubs — мои подписки\n"+
				"• /about — о боте",
			)
		b.api.Send(msg)
	case "/about":
		msg := tgbotapi.NewMessage(chatID, "Я погодный AI агент сюда встроена LLM, поэтому можешь брать советы, допустим по тому как одеться!!!")
		b.api.Send(msg)
	case "/subscribe":
		if len(parts) < 3 {
			msg := tgbotapi.NewMessage(chatID, "❌ Формат: /subscribe [город] [время]\nПример: /subscribe Москва 08:00")
			b.api.Send(msg)
			return
		}

		city := parts[1]
		notifyTime := parts[2]
		if err := b.db.AddSubscription(chatID, city, notifyTime); err != nil {
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка: "+err.Error())
			b.api.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Подписка на %s в %s", city, notifyTime))
		b.api.Send(msg)

		if len(notifyTime) != 5 || notifyTime[2] != ':' {
			msg := tgbotapi.NewMessage(chatID, "❌ Неверный формат времени. Используйте HH:MM (например, 08:00)")
			b.api.Send(msg)
			return
		}

		hour, _ := strconv.Atoi(notifyTime[:2])
		minute, _ := strconv.Atoi(notifyTime[3:])

		if hour < 0 || hour > 23 {
			msg := tgbotapi.NewMessage(chatID, "❌ Часы должны быть от 00 до 23")
			b.api.Send(msg)
			return
		}

		if minute < 0 || minute > 59 {
			msg := tgbotapi.NewMessage(chatID, "❌ Минуты должны быть от 00 до 59")
			b.api.Send(msg)
			return
		}

	case "/unsubscribe":
		if len(parts) < 2 {
			msg := tgbotapi.NewMessage(chatID, "❌ Формат: /unsubscribe [город]")
			b.api.Send(msg)
			return
		}

		city := parts[1]
		if err := b.db.RemoveSubscription(chatID, city); err != nil {
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка: "+err.Error())
			b.api.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Отписка от %s", city))
		b.api.Send(msg)

	case "/mysubs":
		subs, err := b.db.GetUserSubscriptions(chatID)
		log.Printf("📋 Найдено подписок для %d: %d", chatID, len(subs))
		
		if err != nil || len(subs) == 0 {
			msg := tgbotapi.NewMessage(chatID, "Нет активных подписок\nПодписаться: /subscribe Москва 08:00")
			b.api.Send(msg)
			return
		}

		text := "Ваши подписки:\n"
		for _, s := range subs {
			text += fmt.Sprintf("• %s в %s\n", s.City, s.NotifyTime)
		}

		msg := tgbotapi.NewMessage(chatID, text)
		b.api.Send(msg)
		log.Printf("📋 Найдено подписок: %d", len(subs)) 
	}
}
