package handler

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var times int
var activeChats map[string]int64
var lock = sync.RWMutex{}

func writeBotStatistic(l *zap.SugaredLogger, chatTitle string, chatID int64) {
	lock.Lock()
	defer lock.Unlock()
	if activeChats == nil {
		activeChats = make(map[string]int64)
	}
	activeChats[chatTitle] = chatID
	times++
	if times == 200 {
		infoText := "Current active chats:"
		for title, id := range activeChats {
			infoText = infoText + fmt.Sprintf("  Chat title: %s, chat link: %d  ", title, id)
		}
		l.Info(infoText)
		times = 0
	}
}
