.PHONY: build
build:
	@echo ">> building app..."
	go build -v ./cmd/tapmenu

.PHONY: test
test:
	@echo ">> running tests..."
	go mod tidy
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
