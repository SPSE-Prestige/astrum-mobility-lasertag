.PHONY: build run test lint fmt vet tidy docker-build docker-run clean

APP_NAME   := server
BUILD_DIR  := ./cmd/server
BIN_DIR    := ./bin

# ── Build & Run ──

build:
	@echo ">> Building..."
	go build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) $(BUILD_DIR)

run: build
	$(BIN_DIR)/$(APP_NAME)

# ── Quality ──

test:
	go test -race -count=1 -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -1

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

tidy:
	go mod tidy

# ── Docker ──

docker-build:
	docker build -t lasertag-backend .

docker-run: docker-build
	docker run --rm -p 8080:8080 --env-file .env lasertag-backend

# ── Clean ──

clean:
	rm -rf $(BIN_DIR) coverage.out
