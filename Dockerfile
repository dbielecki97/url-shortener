FROM golang:alpine as builder

RUN apk update && apk upgrade && \
apk add --no-cache git
RUN apk add build-base
RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

# Run container
FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8000
CMD ["./server"]