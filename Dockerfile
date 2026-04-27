# Этап 1: Сборка приложения
FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o weather-bot cmd/bot/main.go

FROM golang:1.25-bookworm

WORKDIR /root/

COPY --from=builder /app/weather-bot .

COPY --from=builder /app/migrations ./migrations

#COPY .env .

EXPOSE 8080
CMD ["./weather-bot"]