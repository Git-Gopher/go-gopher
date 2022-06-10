FROM golang:latest AS builder
RUN mkdir -p /workspace
WORKDIR /workspace
COPY . /workspace
RUN ls -la /workspace
RUN go mod download
RUN make build

FROM alpine:latest AS runner
COPY --from=builder /workspace/go-gopher /bin
CMD ["/bin/go-gopher", "analyze"]
ENTRYPOINT ["/entrypoint.sh"]