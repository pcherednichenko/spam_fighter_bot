package app

import (
	"time"

	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/pcherednichenko/spam_fighter_bot/internal/app/data"
	"github.com/pcherednichenko/spam_fighter_bot/internal/app/handler"
)

func StartBot(l *zap.SugaredLogger, token string) {
	// initialize bot
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		l.Fatal("error while initializing bot", err)
	}

	memoryStorage := data.NewMemoryStorage()

	b.Handle(tb.OnUserJoined, handler.UserJoined(l, b, memoryStorage))
	b.Handle(tb.OnText, handler.Text(l, b, memoryStorage))

	b.Handle(tb.OnPhoto, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnAudio, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnAnimation, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnDocument, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnSticker, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnVideo, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnVoice, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnVideoNote, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnContact, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnLocation, handler.Other(l, b, memoryStorage))
	b.Handle(tb.OnVenue, handler.Other(l, b, memoryStorage))

	b.Start()
}
