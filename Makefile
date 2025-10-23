.PHONY: build run test lint clean help

BINARY_NAME=spacdr
MAIN_PATH=main.go

help:
	@echo "Available targets:"
	@echo "  make build    - Build the binary"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make lint     - Run linter"
	@echo "  make clean    - Remove binary"

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	./$(BINARY_NAME)

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY_NAME)
