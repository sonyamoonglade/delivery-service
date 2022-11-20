build-prod:
	docker build -f ./docker/prod.Dockerfile -t sonyamoonglade/sancho-hub:delivery-prod . && docker push sonyamoonglade/sancho-hub:delivery-prod

#cp-env:
#	cp .env.prod ../deployment/delivery/

coverh:
	go tool cover -html=coverage.out && rm coverage.out

.PHONY: cover
cover:
	go test -short -count=1 -coverprofile=coverage.out ./...


gen-mocks:
	mockgen -source=./internal/delivery/storage.go -destination=./internal/delivery/mocks/mock_storage.go && \
	mockgen -source=./internal/runner/storage.go -destination=./internal/runner/mocks/mock_storage.go

run:
	#go build -o ./bin/cli ./cmd/cli/main.go && \
	docker-compose -f ./docker/docker-compose.dev.yaml --env-file ./.env up --build

stop:
	docker-compose -f ./docker/docker-compose.dev.yaml down