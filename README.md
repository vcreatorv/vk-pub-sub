# SubPub: Реализация паттерна Publisher-Subscriber и gRPC-сервиса

## Описание

Данное задание состоит из двух частей:

1. Реализация пакета `subpub`, представляющего собой простую шину событий, основанную на паттерне **Publisher-Subscriber**.
2. Создание gRPC-сервиса подписок, использующего этот пакет.

---

## Часть 1: Пакет `subpub`

### Цель

Реализовать интерфейс шины событий, соответствующий следующему API:

```go
package subpub

import "context"

type MessageHandler func(msg interface{})

type Subscription interface {
	Unsubscribe()
}

type SubPub interface {
	Subscribe(subject string, cb MessageHandler) (Subscription, error)
	Publish(subject string, msg interface{}) error
	Close(ctx context.Context) error
}

func NewSubPub() SubPub {
    panic("Implement me")
}
```

### Требования к реализации
- Один subject может иметь множество подписчиков.

- Медленный подписчик не должен тормозить других.

- Порядок сообщений должен сохраняться (FIFO).

- Метод Close(ctx) должен корректно завершать работу: если контекст отменён — выход немедленно, иначе — дожидаемся завершения публикаций.

- Отсутствие утечек горутин — обязательное условие.

## Часть 2: gRPC-сервис на основе `subpub`
Сервис реализует два метода: подписка на события и публикация события.

Protobuf-схема gRPC сервиса:

```protobuf
syntax = "proto3";
import "google/protobuf/empty.proto";

service PubSub {
  rpc Subscribe(SubscribeRequest) returns (stream Event);
  rpc Publish(PublishRequest) returns (google.protobuf.Empty);
}

message SubscribeRequest {
  string key = 1;
}

message PublishRequest {
  string key = 1;
  string data = 2;
}

message Event {
  string data = 1;
}
```

## Как запускать

Репозиторий: [vk-pub-sub](https://github.com/vcreatorv/vk-pub-sub.git)

### Клонирование репозитория

```bash
git clone https://github.com/vcreatorv/vk-pub-sub.git
cd vk-pub-sub
```

### Установка зависимостей
```bash
go mod tidy
```

### Запуск через Docker

- Сборка и запуск контейнеров:
```bash
make up-build
```

- Запуск без пересборки:
```bash
make up
```

### Запуск вручную (альтернатива)

#### Запуск gRPC-сервиса подписок
```bash
make run-subpub
```

Или
```bash
go run ./cmd/subpub/main.go --config configs/subpub.yml
```

#### Запуск gRPC-клиента
```bash
make run-client
```

Или
```bash
go run ./cmd/app/main.go --config app/app.yml
```

## Тестирование

Для запуска unit-тестов:
```bash
make test
```

Или
```bash
go test -v ./internal/usecase/service
```