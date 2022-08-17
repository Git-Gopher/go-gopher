REPO            := github.com/Git-Gopher/go-gopher
HASH            := $(shell git rev-parse --short HEAD)
DATE            := $(shell date "+%F-%T")
TAG             := $(shell git describe --tags --always --abbrev=0 --match="v[0-9]*.[0-9]*.[0-9]*" 2> /dev/null)
VERSION         := $(shell echo "${TAG}" | sed 's/^.//')
BINARY_DIR      := bin
LDFLAGS_DEV     := -ldflags "-X '${REPO}/version.CommitHash=${HASH}' -X '${REPO}/version.CompileDate=${DATE}'"
LDFLAGS_RELEASE := -ldflags "-s -w -X '${REPO}/version.CommitHash=${HASH}' -X '${REPO}/version.CompileDate=${DATE}'"

.PHONY: run build release install lint format test integration tidy clean

run: build
	go run .

build:
	mkdir -p ${BINARY_DIR}
	go build ${LDFLAGS_DEV} -v -o ${BINARY_DIR}/go-gopher 
	go build ${LDFLAGS_DEV} -v -o ${BINARY_DIR}/go-gopher-cli ./cmd/cli/ 
	go build ${LDFLAGS_DEV} -v -o ${BINARY_DIR}/go-gopher-marker ./cmd/marker/  
	go build ${LDFLAGS_DEV} -v -o ${BINARY_DIR}/go-gopher-github-action ./cmd/github-action/  

release:
	mkdir -p ${BINARY_DIR}
	go build ${LDFLAGS_RELEASE} -v -o ${BINARY_DIR}/go-gopher 
	go build ${LDFLAGS_RELEASE} -v -o ${BINARY_DIR}/go-gopher-cli ./cmd/cli/ 
	go build ${LDFLAGS_RELEASE} -v -o ${BINARY_DIR}/go-gopher-marker ./cmd/marker/  
	go build ${LDFLAGS_RELEASE} -v -o ${BINARY_DIR}/go-gopher-github-action ./cmd/github-action/  

install: build
	cp bin/go-gopher /usr/local/bin

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

clean:
	go clean
	find . -name 'cache.json' -delete
	find . -name 'log-*.json' -delete
	find . -name '*.csv' -delete
	find . -name '*-reports.html' -delete
	rm -rf output || true
