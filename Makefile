# Сборка и запуск Docker-контейнеров
up-build:
	docker-compose up --build

# Запуск без пересборки
up:
	docker-compose up

# Запуск сервиса подписок
run-subpub:
	go run ./cmd/subpub/main.go --config configs/subpub.yml

# Запуск клиента
run-client:
	go run ./cmd/app/main.go --config app/app.yml

# Запуск тестов
test:
	go test -v ./internal/usecase/service