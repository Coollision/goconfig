.PHONY: test cover lint tidy

test:
	go test -race ./...

cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint:
	golangci-lint run

tidy:
	go mod tidy
