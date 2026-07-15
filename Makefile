.PHONY: build run test vet tidy

build:
	go build -o voice-realtime ./cmd/voice-realtime

run:
	go run ./cmd/voice-realtime

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy
