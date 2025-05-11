FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY cmd/subpub/main.go ./cmd/subpub/main.go
COPY configs/subpub.yml ./configs/subpub.yml
COPY internal ./internal

RUN go mod download

RUN go build -o ./.bin ./cmd/subpub/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/configs/subpub.yml ./configs/subpub.yml
COPY --from=builder /app/.bin .

CMD ["./.bin"]