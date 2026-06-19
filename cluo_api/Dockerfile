# syntax=docker/dockerfile:1.4

FROM golang:1.24.5-alpine AS builder
RUN apk update && apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o cluo ./cmd/cluo/*.go && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o seed-admin ./cmd/seed-admin/*.go

# whisper.cpp is invoked via os/exec (internal/infrastructure/whisper/client.go),
# not over the network, so it has to live in this image rather than a sibling
# container. Built from source since no official static binary is published.
FROM alpine:3.20 AS whisper-builder
RUN apk add --no-cache build-base cmake git curl
RUN git clone --depth 1 https://github.com/ggerganov/whisper.cpp.git /whisper.cpp
WORKDIR /whisper.cpp
RUN cmake -B build -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=OFF && \
    cmake --build build --config Release -j"$(nproc)" && \
    find build -maxdepth 3 -type f \( -name 'whisper-cli' -o -name 'main' \) -exec cp {} /usr/local/bin/whisper \; && \
    mkdir -p /models/whisper && \
    curl -fL -o /models/whisper/base.ggml \
      https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin

FROM alpine:latest
RUN apk add --no-cache ca-certificates libstdc++ libgomp

WORKDIR /app

RUN addgroup -g 1000 runner && \
    adduser -D -u 1000 -G runner runner && \
    mkdir -p /app/data && \
    chown -R runner:runner /app

COPY --from=builder /app/cluo .
COPY --from=builder /app/seed-admin .
COPY --from=builder /app/*.env ./
COPY --from=whisper-builder /usr/local/bin/whisper /usr/local/bin/whisper
COPY --from=whisper-builder /models/whisper /models/whisper

USER runner

EXPOSE 8080
CMD ["./cluo"]
