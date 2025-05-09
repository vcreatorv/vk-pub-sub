package main

import (
	"context"
	"github.com/vcreatorv/vk-sub-pub/internal/usecase/service"
	"log"
	"time"
)

const (
	VK       = "VK"
	Telegram = "Telegram"
	Avito    = "Avito"
)

func main() {
	sp := service.NewSubPub()

	sub1, err := sp.Subscribe(VK, func(msg interface{}) {
		log.Printf("%s fast subscriber", VK)
	})
	if err != nil {
		log.Println(err)
		return
	}

	sub2, err := sp.Subscribe(Telegram, func(msg interface{}) {
		log.Printf("%s fast subscriber", Telegram)
	})
	if err != nil {
		log.Println(err)
		return
	}

	sub3, err := sp.Subscribe(Avito, func(msg interface{}) {
		time.Sleep(4 * time.Second)
		log.Printf("%s slow subscriber", Avito)
	})
	if err != nil {
		log.Println(err)
		return
	}

	err = sp.Publish(VK, "Открыт набор на стажировку по Go!")
	if err != nil {
		log.Println(err)
	}

	err = sp.Publish(VK, "Открыт набор на стажировку по Java!")
	if err != nil {
		log.Println(err)
	}

	err = sp.Publish(VK, "Открыт набор на стажировку по ML!")
	if err != nil {
		log.Println(err)
	}

	err = sp.Publish(Telegram, "Открыт набор на стажировку по Go!")
	if err != nil {
		log.Println(err)
	}

	err = sp.Publish(Avito, "Открыт набор на стажировку по Go!")
	if err != nil {
		log.Println(err)
	}

	time.Sleep(2 * time.Second)

	sub1.Unsubscribe()
	sub2.Unsubscribe()
	sub3.Unsubscribe()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	err = sp.Close(ctx)
	if err != nil {
		log.Println(err)
		return
	}
}
