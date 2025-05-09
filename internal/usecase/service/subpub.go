package service

import (
	"context"
	"errors"
	uc "github.com/vcreatorv/vk-sub-pub/internal/usecase"
	"github.com/vcreatorv/vk-sub-pub/internal/utils"
	"log"
	"sync"
)

type SubPubImpl struct {
	topics map[string]map[string]*SubscriptionImpl
	mute   sync.RWMutex
	wg     sync.WaitGroup
	closed bool
}

func NewSubPub() uc.SubPub {
	return &SubPubImpl{
		topics: make(map[string]map[string]*SubscriptionImpl),
	}
}

func (s *SubPubImpl) Subscribe(subject string, cb uc.MessageHandler) (uc.Subscription, error) {
	s.mute.Lock()
	defer s.mute.Unlock()

	if s.closed {
		return nil, errors.New("subpub недоступен")
	}

	id := utils.GenerateID()
	sub := &SubscriptionImpl{
		id:       id,
		subject:  subject,
		handler:  cb,
		messages: make(chan interface{}, 50),
		active:   true,
	}

	if _, ok := s.topics[subject]; !ok {
		s.topics[subject] = make(map[string]*SubscriptionImpl)
	}
	s.topics[subject][id] = sub

	log.Printf("%s пользователь подписался на %s", sub.id, sub.subject)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		sub.Listen()
	}()

	return sub, nil
}

func (s *SubPubImpl) Publish(subject string, msg interface{}) error {
	log.Printf("%s отправил своим подписчикам уведомление %s", subject, msg)

	s.mute.Lock()
	defer s.mute.Unlock()

	if s.closed {
		return errors.New("subpub недоступен")
	}

	subs, ok := s.topics[subject]
	if !ok {
		return errors.New("несуществующий топик")
	}

	for _, sub := range subs {
		sub.mute.Lock()
		select {
		case sub.messages <- msg:
		default:
			log.Printf("%s: подписчик получил максимальное количество уведомлений", sub.id)
		}
		sub.mute.Unlock()
	}

	return nil
}

func (s *SubPubImpl) Close(ctx context.Context) error {
	s.mute.Lock()
	s.closed = true
	s.mute.Unlock()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("все горутины подписок выполнены")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
