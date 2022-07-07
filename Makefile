.PHONY: run
run: build
	go run .

.PHONY: build
build:
	go build -v


.PHONY: lint
lint:
	golangci-lint run

.PHONY: format
format:
	gofumpt -w .

.PHONY: test
test:
	go test -cover ./...

.PHONY: integration
integration:
	go test -cover --tags=integration ./...

.PHONY: tidy
tidy:
	@echo "tidy..."
	go mod tidy
