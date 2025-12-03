# syntax=docker/dockerfile:1.4

FROM golang:1.24.1-alpine AS builder
RUN apk update && apk add --no-cache git upx

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o cluo ./cmd/cluo/*.go && \
    upx --best --lzma cluo # Compress the binary

FROM scratch
WORKDIR /app

COPY --from=builder /app/cluo .
COPY --from=builder /app/migrations/production ./migrations/production/
COPY --from=builder /app/*.env ./

RUN adduser -D runner

RUN mkdir -p /app/data && chown -R runner:runner /app

USER runner

EXPOSE 3500 3501
CMD ["./cluo"]
