include .env

.PHONY: migrate.up migrate.down seed swagger build.core build.seeder build.migrate

migrate.up:
	@go run ./cmd/migrate/main.go up

migrate.down:
	@go run ./cmd/migrate/main.go down

seed:
	@go run ./cmd/seeder/main.go

swagger:
	@swag init -g cmd/core/main.go -o docs/swagger

build.core: swagger
	@go build -o gelora_api cmd/core/main.go

build.seeder:
	@go build -o gelora_seeder cmd/seeder/main.go
