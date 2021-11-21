package handler

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var times int
var activeChats map[int64]string
var lock = sync.RWMutex{}

func writeBotStatistic(l *zap.SugaredLogger, chatID int64, chatTitle string) {
	lock.Lock()
	defer lock.Unlock()
	if activeChats == nil {
		activeChats = make(map[int64]string)
	}
	activeChats[chatID] = chatTitle
	times++
	if times == 200 {
		infoText := "Current active chats:"
		for chatID, chatTitle := range activeChats {
			infoText = infoText + fmt.Sprintf("\nchat id: %d, chat title: %s", chatID, chatTitle)
		}
		l.Infof(infoText)
		times = 0
	}
}
