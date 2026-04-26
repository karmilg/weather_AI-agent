.PHONY: migrate-create migrate-up migrate-down run-cli run-bot build deps docker-up docker-down


run-cli:
	go run cmd/cli/main.go

run-bot:
	go run cmd/bot/main.go

build:
	go build -o bin/weather-cli cmd/cli/main.go
	go build -o bin/weather-bot cmd/bot/main.go

deps:
	go mod tidy
	go mod download

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker compose up --build -d
 
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up
migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down