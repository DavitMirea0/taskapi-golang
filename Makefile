.PHONY: run test lint docker-up docker-down

run:
	go run ./cmd/api

test:
	go test -race ./...

lint:
	golangci-lint run

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v