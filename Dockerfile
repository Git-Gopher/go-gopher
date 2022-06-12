FROM golang:latest AS builder
RUN mkdir -p /workspace
WORKDIR /workspace
COPY . /workspace
RUN go mod download
RUN make build

FROM ubuntu:latest AS runner
RUN mkdir -p /bin
COPY --from=builder /workspace/go-gopher /bin/gitgopher
ENTRYPOINT [ "/bin/gitgopher", "action"]