
build-image:
	go mod tidy
	go mod vendor
	docker build -t server .

build-and-deploy:
	go mod tidy
	go mod vendor
	docker-compose up -d

.PHONY: build-and-deploy
