REPO                  := github.com/Git-Gopher/go-gopher
HASH                  := $(shell git rev-parse --short HEAD)
DATE                  := $(shell date "+%F-%T")
TAG                   := $(shell git describe --tags --always --abbrev=0 --match="v[0-9]*.[0-9]*.[0-9]*" 2> /dev/null)
VERSION               := $(shell echo "${TAG}" | sed 's/^.//')
BINARY_DIR            := bin
LDFLAGS_DEV           := -ldflags "-X '${REPO}/version.CommitHash=${HASH}' -X '${REPO}/version.CompileDate=${DATE}'"
LDFLAGS_RELEASE       := -ldflags "-s -w -X '${REPO}/version.CommitHash=${HASH}' -X '${REPO}/version.CompileDate=${DATE}'"
PACKAGE_NAME          := github.com/Git-Gopher/go-gopher
GOLANG_CROSS_VERSION  ?= v1.19.0


.PHONY: run build release install lint format test integration tidy clean release-windows release-macos release-linux

run: build
	go run ./cmd/cli/

build:
	@mkdir -p ${BINARY_DIR}
ifeq ($(OS),Windows_NT)
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli.exe ./cmd/cli && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker.exe ./cmd/marker && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github.exe ./cmd/github-action && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-workflow.exe ./cmd/workflow
else
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli ./cmd/cli && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker ./cmd/marker && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github-action ./cmd/github-action && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-workflow ./cmd/workflow
endif

release:
	@mkdir -p ${BINARY_DIR}
ifeq ($(OS),Windows_NT)
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli.exe ./cmd/cli && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker.exe ./cmd/marker && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github.exe ./cmd/github-action && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-workflow.exe ./cmd/workflow
else
	@go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli ./cmd/cli && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker ./cmd/marker && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github-action ./cmd/github-action && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-workflow ./cmd/workflow
endif

release-windows:
	@export GOOS=windows && \
	export GOARCH=amd64 && \
	export CGO_ENABLED=1 && \
	echo Building $$GOOS-$$GOARCH && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli-$$GOOS-$$GOARCH.exe ./cmd/cli && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker-$$GOOS-$$GOARCH.exe ./cmd/marker && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github-action-$$GOOS-$$GOARCH.exe ./cmd/github-action && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-workflow-$$GOOS-$$GOARCH.exe ./cmd/workflow

release-macos:
	@export GOOS=darwin && \
	export GOARCH=amd64 && \
	export CGO_ENABLED=1 && \
	echo Building $$GOOS-$$GOARCH && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli-$$GOOS-$$GOARCH.exe ./cmd/cli && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker-$$GOOS-$$GOARCH.exe ./cmd/marker && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github-action-$$GOOS-$$GOARCH.exe ./cmd/github-action && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-workflow-$$GOOS-$$GOARCH.exe ./cmd/workflow


release-linux:
	@export GOOS=linux && \
	export GOARCH=amd64 && \
	export CGO_ENABLED=1 && \
	echo Building $$GOOS-$$GOARCH && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-cli-$$GOOS-$$GOARCH.exe ./cmd/cli && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-marker-$$GOOS-$$GOARCH.exe ./cmd/marker && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-github-action-$$GOOS-$$GOARCH.exe ./cmd/github-action && \
	go build ${LDFLAGS_RELEASE} -o ${BINARY_DIR}/go-gopher-workflow-$$GOOS-$$GOARCH.exe ./cmd/workflow


install: build
	cp bin/go-gopher ~/.local/bin

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


.PHONY: go-releaser-release-dry-run
go-releaser-release-dry-run:
	@docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate --skip-publish

.PHONY: go-releaser-release
go-releaser-release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist