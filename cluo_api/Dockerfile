# syntax=docker/dockerfile:1.4

FROM golang:1.24.5-alpine AS builder
RUN apk update && apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o cluo ./cmd/cluo/*.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app

RUN addgroup -g 1000 runner && \
    adduser -D -u 1000 -G runner runner && \
    mkdir -p /app/data && \
    chown -R runner:runner /app

COPY --from=builder /app/cluo .
COPY --from=builder /app/*.env ./

USER runner

EXPOSE 8080
CMD ["./cluo"]
