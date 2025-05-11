package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/vcreatorv/vk-sub-pub/internal/config"
	"github.com/vcreatorv/vk-sub-pub/internal/delivery/grpc/interceptors"
	pb "github.com/vcreatorv/vk-sub-pub/internal/delivery/grpc/subpub/proto"
	"github.com/vcreatorv/vk-sub-pub/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

const (
	VK       = "VK"
	Telegram = "Telegram"
	Avito    = "Avito"
)

var configPath = flag.String("config", "configs/app.yml", "путь до файла конфигураций")

func LoadConfig() (*config.AppConfig, error) {
	flag.Parse()

	yamlFile, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	var cfg config.AppConfig
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга yaml файла: %w", err)
	}

	return &cfg, nil
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	grpcConn, err := grpc.NewClient(
		cfg.GetAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(interceptors.TimeoutClientInterceptor()),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer grpcConn.Close()

	subPubClient := pb.NewPubSubClient(grpcConn)

	go func() {
		ctx := context.Background()
		ctx = utils.SetTimeout(ctx, 2*time.Second)
		stream, err := subPubClient.Subscribe(ctx, &pb.SubscribeRequest{Key: VK})
		if err != nil {
			log.Println(err)
			return
		}

		for {
			event, err := stream.Recv()
			if err != nil {
				log.Printf("Stream error: %v", err)
				return
			}
			log.Printf("%s Пришло сообщение: %s", VK, event.Data)
		}
	}()

	for i := 0; i < 35; i++ {
		_, err := subPubClient.Publish(context.Background(), &pb.PublishRequest{
			Key:  VK,
			Data: fmt.Sprintf("Стажировка #%d", i+1),
		})
		if err != nil {
			log.Print(err)
		}
	}
	time.Sleep(10 * time.Second)
	log.Println("Завершение работы клиента...")
}
