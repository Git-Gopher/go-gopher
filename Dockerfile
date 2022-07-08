FROM golang:latest AS builder
RUN mkdir -p /workspace
WORKDIR /workspace
COPY . /workspace
RUN go mod download
RUN make build
FROM ubuntu:latest AS runner
RUN apt-get update \
    && apt-get -y upgrade \
    && apt-get install -y \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true
RUN mkdir -p /bin
COPY --from=builder /workspace/go-gopher /bin/gitgopher
ENTRYPOINT [ "/bin/gitgopher", "action", "--logging"]