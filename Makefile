#!make
include .env
export


build-app:
	go build -o ./bin/app.exe ./cmd/app/main.go

build-cli:
	go build -o ./bin/cli.exe ./cmd/cli/main.go

run:
	go run cmd/app/main.go

build-delivery-local:
	docker build -f ./docker/local.Dockerfile -t sonyamoonglade/sancho-hub:delivery-local .

build-delivery-prod:
	docker build -f ./docker/prod.Dockerfile -t sonyamoonglade/sancho-hub:delivery-prod .

run-delivery-local:
	docker run -d -p 9000:9000 --env-file ./.env.local sonyamoonglade/sancho-hub:delivery-local

run-delivery-prod:
	docker run -d -p 9000:9000 --env-file ./.env.prod sonyamoonglade/sancho-hub:delivery-prod

cp-env:
	cp .env.prod ../../sancho-console/delivery/
