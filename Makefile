REPO            := github.com/Git-Gopher/go-gopher
HASH            := $(shell git rev-parse --short HEAD)
DATE            := $(shell date "+%F-%T")
TAG             := $(shell git describe --tags --always --abbrev=0 --match="v[0-9]*.[0-9]*.[0-9]*" 2> /dev/null)
VERSION         := $(shell echo "${TAG}" | sed 's/^.//')
BINARY_DIR      := bin
LDFLAGS_DEV     := -ldflags "-X '${REPO}/version.CommitHash=${HASH}' -X '${REPO}/version.CompileDate=${DATE}'"
LDFLAGS_RELEASE := -ldflags "-s -w -X '${REPO}/version.CommitHash=${HASH}' -X '${REPO}/version.CompileDate=${DATE}'"

.PHONY: run build release install lint format test integration tidy clean release-windows release-macos release-linux

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
	$(MAKE) release-linux
	$(MAKE) release-macos
	$(MAKE) release-windows

release-windows:
	@GOOS=windows
	@GOARCH=amd64

	echo "Building $GOOS-$GOARCH"

	@GOOS=${GOOS} GOARCH=${GOARCH} go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli-${GOOS}-${GOARCH}.exe ./cmd/cli
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker-${GOOS}-${GOARCH}.exe ./cmd/marker
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github-action-${GOOS}-${GOARCH}.exe ./cmd/github-action

release-macos:
	export GOOS=darwin
	export GOARCH=amd64

	@echo Building $$GOOS-$$GOARCH

	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli-${GOOS}-${GOARCH} ./cmd/cli
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker-${GOOS}-${GOARCH} ./cmd/marker
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github-action-${GOOS}-${GOARCH} ./cmd/github-action

release-linux:
	export GOOS=linux
	export GOARCH=amd64

	@echo Building $$GOOS-$$GOARCH

	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli-${GOOS}-${GOARCH} ./cmd/cli
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker-${GOOS}-${GOARCH} ./cmd/marker
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github-action-${GOOS}-${GOARCH} ./cmd/github-action


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
	find . -name 'log-*.json' -delete
	find . -name '*.csv' -delete
	rm -rf output
