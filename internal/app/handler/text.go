package handler

import (
	"strconv"
	"time"

	"github.com/pcherednichenko/spam_fighter_bot/internal/app/data"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Text(l *zap.SugaredLogger, b *tb.Bot, s data.Storage) func(m *tb.Message) {
	return func(m *tb.Message) {
		info, ok := s.Exist(m.Chat, m.Sender)
		if !ok {
			return
		}
		if m.Text != strconv.Itoa(info.RightAnswer) {
			err := b.Delete(m)
			if err != nil {
				l.Errorf("error while deleting (spam) user message: %v", err)
			}
			return
		}
		// in case of correct answer:
		s.Remove(m.Chat, m.Sender)
		// Correct! Tell us about yourself
		approveMessage, err := b.Send(m.Chat, "Верно! Расскажите нам о себе 🙂")
		if err != nil {
			l.Errorf("error while sending: %v", err)
		}
		go deleteWelcomeMessages(l, b, m, approveMessage, info.WelcomeMessage)
	}
}

func deleteWelcomeMessages(l *zap.SugaredLogger, b *tb.Bot, m *tb.Message, approveMessage *tb.Message, welcomeMessage *tb.Message) {
	time.Sleep(time.Second * 30)
	err := b.Delete(m)
	if err != nil {
		l.Errorf("error while deleting user message: %v", err)
	}
	err = b.Delete(approveMessage)
	if err != nil {
		l.Errorf("error while deleting approve message: %v", err)
	}
	err = b.Delete(welcomeMessage)
	if err != nil {
		l.Errorf("error while deleting welcome message after approve: %v", err)
	}
}
