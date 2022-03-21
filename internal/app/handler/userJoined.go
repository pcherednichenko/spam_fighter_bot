package handler

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	ntw "moul.io/number-to-words"

	"github.com/pcherednichenko/spam_fighter_bot/internal/app/data"
)

const waitForAnAnswerTime = time.Minute * 3

func UserJoined(l *zap.SugaredLogger, b *tb.Bot, s data.Storage) func(m *tb.Message) {
	return func(m *tb.Message) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		// generating two random small number
		firstNumber := r.Intn(4) + 1
		secondNumber := r.Intn(4) + 1
		fistNumberInWordsRu := ntw.IntegerToRuRu(firstNumber)
		fistNumberInWordsEn := ntw.IntegerToEnUs(firstNumber)
		secondNumberInWordsRu := ntw.IntegerToRuRu(secondNumber)
		secondNumberInWordsEn := ntw.IntegerToEnUs(secondNumber)

		if excludedUser(m.UserJoined) {
			return
		}
		username := getUsername(m.UserJoined)
		welcomeMessageText := getWelcomeMessageText(username, m.Chat.Title,
			fistNumberInWordsEn, secondNumberInWordsEn, fistNumberInWordsRu, secondNumberInWordsRu)
		welcomeMessage, err := b.Send(m.Chat, welcomeMessageText)
		if err != nil {
			l.Error("error while sending welcome message", err)
			return
		}
		s.Add(m.Chat, m.UserJoined, data.Info{WelcomeMessage: welcomeMessage, RightAnswer: firstNumber * secondNumber})
		// Goroutine to delete message after waitForAnAnswerTime minutes
		// and block user if he or she still in the list
		go checkAndBanUser(l, b, welcomeMessage, s, m, username)
	}
}

func excludedUser(u *tb.User) bool {
	// we trust these bots
	if u.Username == "shieldy_bot" {
		return true
	}
	if u.Username == "combot" {
		return true
	}
	return false
}

func checkAndBanUser(l *zap.SugaredLogger, b *tb.Bot, welcomeMessage *tb.Message, s data.Storage, m *tb.Message, username string) {
	time.Sleep(waitForAnAnswerTime)
	err := b.Delete(welcomeMessage)
	if err != nil {
		// maybe message already deleted because of correct answer
		if !strings.Contains(err.Error(), "message to delete not found") {
			l.Errorf("error while deleting welcome message after time: %v", err)
		}
	}
	if _, ok := s.Exist(m.Chat, m.UserJoined); ok {
		userToBan, err := b.ChatMemberOf(m.Chat, m.UserJoined)
		if err != nil {
			l.Errorf("error while banning user, chat title: %s, error: %v", m.Chat.Title, err)
			return
		}
		if userToBan == nil {
			l.Errorf("error while restricting user, chat title: %s, error: %v", m.Chat.Title, err)
			return
		}
		err = b.Restrict(m.Chat, userToBan)
		if err != nil {
			l.Errorf("error while restricting user, chat title: %s, error: %v", m.Chat.Title, err)
		}
		err = b.Ban(m.Chat, userToBan)
		if err != nil {
			l.Errorf("error while ban, chat title: %s, error: %v", m.Chat.Title, err)
		}
		err = b.Delete(m)
		if err != nil {
			l.Errorf("error while deleting joining message, chat title: %s, error: %v", m.Chat.Title, err)
		}
		l.Infof("Banned: %s", username)
	}
}

func getUsername(u *tb.User) string {
	if u.Username != "" {
		return "@" + u.Username
	}
	username := ""
	username = u.FirstName
	if u.LastName != "" {
		username = username + " " + u.LastName
	}
	return username
}

func getWelcomeMessageText(username, chatName,
	fistNumberInWordsEn, secondNumberInWordsEn, fistNumberInWordsRu, secondNumberInWordsRu string) string {
	// Welcome to the chat of Russian-speaking people in Amsterdam!
	// Write the number as: %s multiply by %s to check that you are not a bot
	firstNumText, secondNumText := fistNumberInWordsEn, secondNumberInWordsEn
	welcomeText := "Hi %s, welcome to the %s! " +
		"To make sure that you aren't a bot, please answer " +
		"the following question by replying with the result as a number: %s multiplied by %s"
	if chatNameContainsCyrillic(chatName) {
		firstNumText, secondNumText = fistNumberInWordsRu, secondNumberInWordsRu
		welcomeText = "%s Добро пожаловать в чат %s! " +
			"Напишите числом сколько будет: %s умножить на %s, чтобы проверить, что вы не бот"
	}
	return fmt.Sprintf(
		welcomeText,
		username, chatName, firstNumText, secondNumText)
}

func chatNameContainsCyrillic(chatName string) bool {
	for _, char := range chatName {
		if unicode.Is(unicode.Cyrillic, char) {
			return true
		}
	}
	return false
}
