package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/vcreatorv/vk-sub-pub/internal/config"
	"github.com/vcreatorv/vk-sub-pub/internal/delivery/grpc/interceptors"
	"github.com/vcreatorv/vk-sub-pub/internal/delivery/grpc/subpub"
	pb "github.com/vcreatorv/vk-sub-pub/internal/delivery/grpc/subpub/proto"
	"github.com/vcreatorv/vk-sub-pub/internal/usecase/service"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath = flag.String("config", "configs/subpub.yml", "путь до файла конфигураций")

func LoadConfig() (*config.SubPubConfig, error) {
	flag.Parse()

	yamlFile, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	var cfg config.SubPubConfig
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

	subPubService := service.NewSubPub()

	subPubGRPC := subpub.NewGRPC(subPubService)

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(interceptors.TimeoutServerInterceptor()),
	)
	pb.RegisterPubSubServer(grpcServer, subPubGRPC)

	lis, err := net.Listen("tcp", cfg.GetAddr())
	if err != nil {
		log.Fatalf("Невозможно прослушать порт: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Ошибка сервера gRPC:", err)
		}
	}()

	log.Printf("Запуск сервиса подписок по адресу %s", cfg.GetAddr())

	sig := <-quit
	log.Printf("Завершение работы сервиса подписок по сигналу: %v", sig)

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := subPubService.Close(ctx); err != nil {
		log.Printf("Предупреждение при досрочном закрытии сервиса: %v (хендлеры продолжают работу)", err)
	}

	<-stopped
	log.Println("Сервис подписок полностью остановлен")
}
