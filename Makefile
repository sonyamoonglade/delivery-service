#!make
include .env
export


build-app:
	go build -o ./bin/app.exe ./cmd/app/main.go

build-cli:
	go build -o ./bin/cli.exe ./cmd/cli/main.go

run:
	go run cmd/app/main.go

build-local:
	docker build -f ./docker/local.Dockerfile -t sonyamoonglade/sancho-hub:delivery-local . && docker push sonyamoonglade/sancho-hub:delivery-local

build-prod:
	docker build -f ./docker/prod.Dockerfile -t sonyamoonglade/sancho-hub:delivery-prod . && docker push sonyamoonglade/sancho-hub:delivery-prod

run-local:
	docker run -d -p 9000:9000 --env-file ./.env.local sonyamoonglade/sancho-hub:delivery-local

run-prod:
	docker run -d -p 9000:9000 --env-file ./.env.prod sonyamoonglade/sancho-hub:delivery-prod

cp-env:
	cp .env.prod ../deployment/delivery/

coverh:
	go tool cover -html=coverage.out && rm coverage.out

.PHONY: cover
cover:
	go test -short -count=1 -coverprofile=coverage.out ./...
