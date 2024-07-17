# Targets
.PHONY: all build run docker-build docker-run docker-stop migrate-up migrate-down docker-logs docker-restart

all: build run

build:
	go build -o main ./cmd/api

run:
	./main
	
test:
	go test -v ./...

docker-build:
	docker build -t my-go-app .

docker-run:
	docker-compose up

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-restart:
	docker-compose restart

migrate-up:
	docker-compose run --rm app migrate -path /migrations -database "postgres://$(DB_USER):$(DB_PASS)@db:5432/$(DB_NAME)?sslmode=disable" up

migrate-down:
	docker-compose run --rm app migrate -path /migrations -database "postgres://$(DB_USER):$(DB_PASS)@db:5432/$(DB_NAME)?sslmode=disable" down
