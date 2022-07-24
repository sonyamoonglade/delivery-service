#!make
include .env
export


run:
	go run cmd/app/main.go

migrate-up:
	migrate -database "$(DB_URL)" -path migrations -verbose up

migrate-down:
	migrate -database "$(DB_URL)" -path migrations -verbose down

run-db:
	docker-compose -f docker/docker-compose.yaml up --detach

stop-db:
	docker-compose -f docker/docker-compose.yaml down

local-migrate-up:
	migrate -path migrations -database "postgresql://dlvuser:dlvuserpwd@localhost:5433/delivery?sslmode=disable" -verbose up

local-migrate-down:
	migrate -database "postgresql://dlvuser:dlvuserpwd@localhost:5433/delivery?sslmode=disable" -path migrations -verbose down


