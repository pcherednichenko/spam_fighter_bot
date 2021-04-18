package data

import (
	"fmt"
	"sync"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Storage interface {
	Add(chat *tb.Chat, user *tb.User, i Info)
	Remove(chat *tb.Chat, user *tb.User)
	Exist(chat *tb.Chat, user *tb.User) (Info, bool)
}

func NewMemoryStorage() *MemoryStorage {
	m := MemoryStorage{}
	m.pending = make(map[string]Info)
	return &m
}

// MemoryStorage storing in memory for simplicity
type MemoryStorage struct {
	sync.RWMutex
	pending map[string]Info
}

type Info struct {
	WelcomeMessage *tb.Message
	RightAnswer    int
}

func (s *MemoryStorage) Add(chat *tb.Chat, user *tb.User, i Info) {
	s.Lock()
	defer s.Unlock()
	s.pending[key(chat, user)] = i
}

func (s *MemoryStorage) Remove(chat *tb.Chat, user *tb.User) {
	s.Lock()
	defer s.Unlock()
	delete(s.pending, key(chat, user))
}

func (s *MemoryStorage) Exist(chat *tb.Chat, user *tb.User) (Info, bool) {
	s.Lock()
	defer s.Unlock()
	info, ok := s.pending[key(chat, user)]
	return info, ok
}

func key(chat *tb.Chat, u *tb.User) string {
	return fmt.Sprintf("%d-%d", chat.ID, u.ID)
}
