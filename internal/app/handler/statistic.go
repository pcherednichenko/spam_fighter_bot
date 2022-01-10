package handler

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var times int
var activeChats map[int64]string
var lock = sync.RWMutex{}

func writeBotStatistic(l *zap.SugaredLogger, chatTitle string, chatID int64) {
	lock.Lock()
	defer lock.Unlock()
	if activeChats == nil {
		activeChats = make(map[int64]string)
	}
	activeChats[chatID] = chatTitle
	times++
	if times == 200 {
		infoText := "Current active chats:"
		for id, title := range activeChats {
			infoText = infoText + fmt.Sprintf("  Chat id: %d, chat title: %s  ", id, title)
		}
		l.Info(infoText)
		times = 0
	}
}
