.PHONY: run
run: build
	go run .

.PHONY: build
build:
	mkdir bin || true
	go build -v -o bin/go-gopher
	go build -v -o bin/go-gopher-cli ./cmd/cli/ 
	go build -v -o bin/go-gopher-marker ./cmd/marker/  
	go build -v -o bin/go-gopher-github-action ./cmd/github-action/  

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

.PHONY: clean
clean:
	find . -name 'log-*.json' -delete
	find . -name '*.csv' -delete
	rm -rf output
