REPO            := github.com/Git-Gopher/go-gopher
HASH            := $(shell git rev-parse --short HEAD)
DATE            := $(shell date "+%F-%T")
TAG             := $(shell git describe --tags --always --abbrev=0 --match="v[0-9]*.[0-9]*.[0-9]*" 2> /dev/null)
VERSION         := $(shell echo "${TAG}" | sed 's/^.//')
BINARY_DIR      := bin
LDFLAGS_DEV     := -ldflags "-X '${REPO}/version.CommitHash=${HASH}' -X '${REPO}/version.CompileDate=${DATE}'"
LDFLAGS_RELEASE := -ldflags "-s -w -X '${REPO}/version.CommitHash=${HASH}' -X '${REPO}/version.CompileDate=${DATE}'"

.PHONY: run
run: build
	go run .

.PHONY: build
build:
	mkdir ${BINARY_DIR} || true
	go build ${LDFLAGS_DEV} -v -o ${BINARY_DIR}/go-gopher 
	go build ${LDFLAGS_DEV} -v -o ${BINARY_DIR}/go-gopher-cli ./cmd/cli/ 
	go build ${LDFLAGS_DEV} -v -o ${BINARY_DIR}/go-gopher-marker ./cmd/marker/  
	go build ${LDFLAGS_DEV} -v -o ${BINARY_DIR}/go-gopher-github-action ./cmd/github-action/  

.PHONY: release 
release:
	mkdir ${BINARY_DIR} || true
	go build ${LDFLAGS_RELEASE} -v -o ${BINARY_DIR}/go-gopher 
	go build ${LDFLAGS_RELEASE} -v -o ${BINARY_DIR}/go-gopher-cli ./cmd/cli/ 
	go build ${LDFLAGS_RELEASE} -v -o ${BINARY_DIR}/go-gopher-marker ./cmd/marker/  
	go build ${LDFLAGS_RELEASE} -v -o ${BINARY_DIR}/go-gopher-github-action ./cmd/github-action/  

.PHONY: install
install: build
	cp bin/go-gopher /usr/local/bin

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
	go clean
	find . -name 'log-*.json' -delete
	find . -name '*.csv' -delete
	rm -rf output
