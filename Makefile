.PHONY: build run test clean tidy

BINARY_NAME=vendor
LDFLAGS="-s -w"

build:
	CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o bin/$(BINARY_NAME) cmd/vendor/main.go

run:
	go run cmd/vendor/main.go

test:
	go test ./...

clean:
	go clean
	rm -rf bin/

tidy:
	go mod tidy

