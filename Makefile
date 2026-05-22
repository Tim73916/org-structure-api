.PHONY: run docker-up docker-down test build clean

run:
	@go run cmd/main.go

build:
	@go build -o bin/app cmd/main.go

docker-up:
	@docker-compose up -d --build

docker-down:
	@docker-compose down

test:
	@go test -v ./...

clean:
	@rm -rf bin/