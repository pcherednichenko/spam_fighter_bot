package handler

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var times int
var activeChats map[string]string
var lock = sync.RWMutex{}

func writeBotStatistic(l *zap.SugaredLogger, chatTitle, link string) {
	lock.Lock()
	defer lock.Unlock()
	if activeChats == nil {
		activeChats = make(map[string]string)
	}
	activeChats[chatTitle] = link
	times++
	if times == 200 {
		infoText := "Current active chats:"
		for title, link := range activeChats {
			infoText = infoText + fmt.Sprintf("  Chat title: %s, chat link: %s  ", title, link)
		}
		l.Info(infoText)
		times = 0
	}
}
