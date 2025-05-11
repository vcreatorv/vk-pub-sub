package service

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSubscription_Unsubscribe(t *testing.T) {
	t.Parallel()

	sub := &SubscriptionImpl{
		id:       "12DF17ER",
		subject:  "VK",
		handler:  func(msg interface{}) {},
		messages: make(chan interface{}),
		active:   true,
	}

	sub.Unsubscribe()

	// Повторный вызов
	require.NotPanics(t, func() {
		sub.Unsubscribe()
	})

	require.False(t, sub.active, "подписка должна быть неактивной после отписки")
}
