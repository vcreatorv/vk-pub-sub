package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestSubPub_Subscribe(t *testing.T) {
	t.Parallel()

	t.Run("Успешная подписка", func(t *testing.T) {
		t.Parallel()
		sp := NewSubPub()
		sub, err := sp.Subscribe("VK", func(msg interface{}) {})
		require.NoError(t, err)
		require.NotNil(t, sub)
	})

	t.Run("Подписка при недоступном сервисе подписок", func(t *testing.T) {
		t.Parallel()
		sp := NewSubPub()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := sp.Close(ctx)
		require.NoError(t, err)

		_, err = sp.Subscribe("VK", func(msg interface{}) {})
		require.Error(t, err)
		require.Equal(t, "сервис подписок недоступен", err.Error())
	})

	t.Run("Все подписчики получили сообщение", func(t *testing.T) {
		t.Parallel()

		sp := NewSubPub()
		var subscriberCount = 5
		var wg sync.WaitGroup
		received := make([]bool, subscriberCount)

		for i := 0; i < subscriberCount; i++ {
			wg.Add(1)
			idx := i
			_, err := sp.Subscribe("VK", func(msg interface{}) {
				defer wg.Done()
				received[idx] = true
			})
			require.NoError(t, err)
		}

		err := sp.Publish("VK", "Открылась стажировка по Go!")
		require.NoError(t, err)

		wg.Wait()
		for i := 0; i < subscriberCount; i++ {
			require.True(t, received[i])
		}
	})
}

func TestSubPub_Publish(t *testing.T) {
	t.Parallel()

	t.Run("Успешная публикация", func(t *testing.T) {
		t.Parallel()
		sp := NewSubPub()
		var wg sync.WaitGroup
		wg.Add(1)

		_, err := sp.Subscribe("VK", func(msg interface{}) {
			defer wg.Done()
			require.Equal(t, "Открылась стажировка по Go!", msg)
		})
		require.NoError(t, err)

		err = sp.Publish("VK", "Открылась стажировка по Go!")
		require.NoError(t, err)

		wg.Wait()
	})

	t.Run("Публикация в несуществующий топик", func(t *testing.T) {
		t.Parallel()
		sp := NewSubPub()
		err := sp.Publish("ABC", "Открылась стажировка по Go!")
		require.Error(t, err)
		require.Equal(t, "несуществующий топик", err.Error())
	})

	t.Run("Публикация при недоступной сервисе подписок", func(t *testing.T) {
		t.Parallel()
		sp := NewSubPub()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := sp.Close(ctx)
		require.NoError(t, err)

		err = sp.Publish("VK", "Открылась стажировка по Go!")
		require.Error(t, err)
		require.Equal(t, "сервис подписок недоступен", err.Error())
	})
}

func TestSubPub_Close(t *testing.T) {
	t.Parallel()

	t.Run("Успешное закрытие", func(t *testing.T) {
		t.Parallel()
		sp := NewSubPub()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := sp.Close(ctx)
		require.NoError(t, err)
	})

	t.Run("Закрытие с таймаутом", func(t *testing.T) {
		t.Parallel()
		sp := NewSubPub()

		_, err := sp.Subscribe("VK", func(msg interface{}) {})
		require.NoError(t, err)

		err = sp.Publish("VK", "Открылась стажировка по Go!")
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err = sp.Close(ctx)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
