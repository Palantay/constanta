
.PHONY: go
go: build start

.PHONY: api-build
api-build: dc-up migrate-up build start

.PHONY: build
build:
	go build -v ./cmd/api/

.PHONY: migrate up
migrate-up:
	migrate -path migrations -database "postgres://localhost:5432/db?sslmode=disable&user=admin&password=admin" up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "postgres://localhost:5432/db?sslmode=disable&user=admin&password=admin" down

.PHONY: dc-up
dc-up:
	docker-compose up -d

.PHONY: dc-down
dc-down:
	docker-compose down

.PHONY: start
start: 
	./api

.DEFAULT_GOAL := go