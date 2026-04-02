# ── Stage 1: Build ──
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# ── Stage 2: Runtime ──
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=builder /app/server .

USER app
EXPOSE 8080
ENTRYPOINT ["./server"]
