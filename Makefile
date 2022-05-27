run:
	go run .

build:
	go build

lint:
	golangci-lint run

format:
	gofumpt -w .

test:
	go test -cover ./...

integration:
	go test -cover --tags=integration ./...

tidy:
	@echo "tidy..."
	go mod tidy