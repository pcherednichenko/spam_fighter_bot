package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/pcherednichenko/spam_fighter_bot/internal/app/data"
)

func Text(l *zap.SugaredLogger, b *tb.Bot, s data.Storage) func(m *tb.Message) {
	return func(m *tb.Message) {
		if m.Private() {
			err := sendPrivateMessageResponse(b, m)
			if err != nil {
				l.Errorf("error while sending private message response: %v", err)
			}
			return
		}
		if m.Chat != nil {
			writeBotStatistic(l, m.Chat.Title, m.Chat.ID)
		}
		info, ok := s.Exist(m.Chat, m.Sender)
		if !ok {
			return
		}
		if strings.TrimSpace(m.Text) != strconv.Itoa(info.RightAnswer) {
			err := b.Delete(m)
			if err != nil {
				l.Errorf("error while deleting (spam) user message: %v", err)
			}
			return
		}
		// in case of correct answer:
		s.Remove(m.Chat, m.Sender)
		// Correct! Tell us about yourself
		correctText := "Correct!"
		tellUsAboutYourselfText := "%s tell us about yourself 🙂"
		if chatNameContainsCyrillic(m.Chat.Title) {
			correctText = "Верно!"
			tellUsAboutYourselfText = "%s расскажите нам о себе 🙂"
		}
		approveMessage, err := b.Send(m.Chat, correctText)
		if err != nil {
			l.Errorf("error while sending: %v", err)
		}
		// imitation of real typing delays
		time.Sleep(time.Second * 2)
		tellUsText := fmt.Sprintf(tellUsAboutYourselfText, getUsername(m.Sender))
		if strings.Contains(m.Chat.Title, "Амстердам") {
			// Additional text for chat tell us message
			tellUsText += " Чем занимаетесь в Амстердаме?"
		}
		tellUsMessage, err := b.Send(m.Chat, tellUsText)
		if err != nil {
			l.Errorf("error while sending: %v", err)
		}
		go deleteWelcomeMessages(l, b, m, approveMessage, info.WelcomeMessage, tellUsMessage)
	}
}

func deleteWelcomeMessages(l *zap.SugaredLogger, b *tb.Bot,
	m *tb.Message, approveMessage *tb.Message, welcomeMessage *tb.Message, tellUsMessage *tb.Message,
) {
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
	// delay before deleting second welcome message
	time.Sleep(time.Second * 90)
	if tellUsMessage == nil {
		l.Errorf("error while deleting tell us message: %v", err)
		return
	}
	err = b.Delete(tellUsMessage)
	if err != nil {
		l.Errorf("error while deleting tell us about yourself message after approve: %v", err)
	}
}

func sendPrivateMessageResponse(b *tb.Bot, m *tb.Message) error {
	_, err := b.Send(m.Chat, "Hello! Just add me to the chat, make me an admin"+
		" and grant the rights to delete messages and ban spam users")
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 3)
	_, err = b.Send(m.Chat, "I hope you enjoy my work! 😉")
	if err != nil {
		return err
	}
	return nil
}
