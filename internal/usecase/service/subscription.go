package service

import (
	uc "github.com/vcreatorv/vk-sub-pub/internal/usecase"
	"log"
	"sync"
)

type SubscriptionImpl struct {
	id       string
	subject  string
	handler  uc.MessageHandler
	messages chan interface{}
	active   bool
	mute     sync.RWMutex
}

func (s *SubscriptionImpl) Unsubscribe() {
	s.mute.Lock()
	defer s.mute.Unlock()
	if s.active {
		s.active = false
		log.Printf("%s подписчик отписался от %s", s.id, s.subject)
		close(s.messages)
	}
}

func (s *SubscriptionImpl) Listen() {
	for msg := range s.messages {
		if !s.active {
			return
		}
		log.Printf("%s подписчик получил уведомление от %s: %v", s.id, s.subject, msg)
		s.handler(msg)
	}
}
