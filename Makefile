.PHONY: tidy test build run-product run-category run-attribute docker-up docker-down

tidy:
	go mod tidy

test:
	go test ./...

build:
	go build -o bin/product ./cmd/product
	go build -o bin/category ./cmd/category
	go build -o bin/attribute ./cmd/attribute

run-product:
	go run ./cmd/product

run-category:
	go run ./cmd/category

run-attribute:
	go run ./cmd/attribute

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down
