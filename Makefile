test:
	go test ./... -v -race

build:
	go build cmd/server.go
build-image:
	go mod tidy
	go mod vendor
	docker build -t server .

build-and-deploy:
	go mod tidy
	go mod vendor
	docker-compose up -d

.PHONY: test build
