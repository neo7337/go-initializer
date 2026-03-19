# ── go-initializer Makefile ────────────────────────────────────────────────────
# Usage: make <target>
# Run `make help` for a list of available targets.

.DEFAULT_GOAL := help
BINARY_SERVER := bin/server
BINARY_CLI    := bin/goini

# ── Help ───────────────────────────────────────────────────────────────────────

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

# ── Go — dependencies ─────────────────────────────────────────────────────────

.PHONY: tidy
tidy: ## Tidy and verify Go module dependencies
	go mod tidy
	go mod verify

# ── Go — build ────────────────────────────────────────────────────────────────

.PHONY: build
build: build-server build-cli ## Build both the server and goini binaries

.PHONY: build-server
build-server: ## Build the HTTP server binary → bin/server
	go build -o $(BINARY_SERVER) ./cmd/server/...

.PHONY: build-cli
build-cli: ## Build the goini CLI binary → bin/goini
	go build -o $(BINARY_CLI) ./cmd/goini/...

# ── Go — run ──────────────────────────────────────────────────────────────────

.PHONY: run-server
run-server: ## Run the HTTP server (hot-reload not included; uses go run)
	go run ./cmd/server/...

.PHONY: run-cli
run-cli: ## Run the goini CLI (pass ARGS="new --name foo" to forward flags)
	go run ./cmd/goini/... $(ARGS)

# ── Go — test ─────────────────────────────────────────────────────────────────

.PHONY: test
test: ## Run all Go tests with race detector
	go test -race ./...

.PHONY: test-verbose
test-verbose: ## Run all Go tests with verbose output
	go test -race -v ./...

.PHONY: test-cover
test-cover: ## Run tests and show coverage report
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# ── Go — lint / vet ───────────────────────────────────────────────────────────

.PHONY: vet
vet: ## Run go vet on all packages
	go vet ./...

# ── Frontend ──────────────────────────────────────────────────────────────────

.PHONY: frontend-install
frontend-install: ## Install frontend npm dependencies
	cd frontend && npm install

.PHONY: frontend-dev
frontend-dev: ## Start frontend Vite dev server (http://localhost:3000)
	cd frontend && npm run start

.PHONY: frontend-build
frontend-build: ## Build frontend for production → frontend/dist/
	cd frontend && npm run build

.PHONY: frontend-preview
frontend-preview: ## Preview the production frontend build locally
	cd frontend && npm run preview

# ── Full stack (local, no Docker) ─────────────────────────────────────────────

.PHONY: dev
dev: ## Start backend server + frontend dev server concurrently
	@trap 'kill 0' SIGINT; \
	$(MAKE) run-server & \
	$(MAKE) frontend-dev & \
	wait

# ── Docker ────────────────────────────────────────────────────────────────────

.PHONY: docker-build
docker-build: ## Build all Docker images via docker-compose
	docker compose build

.PHONY: docker-up
docker-up: ## Start all services via docker-compose (detached)
	docker compose up -d

.PHONY: docker-down
docker-down: ## Stop and remove all docker-compose containers
	docker compose down

.PHONY: docker-logs
docker-logs: ## Tail logs from all docker-compose services
	docker compose logs -f

.PHONY: docker-restart
docker-restart: docker-down docker-up ## Rebuild images and restart all services

# ── goini CLI helpers ─────────────────────────────────────────────────────────

.PHONY: install-cli
install-cli: ## Install goini CLI to GOPATH/bin via go install
	go install ./cmd/goini/...

.PHONY: completions
completions: build-cli ## Generate shell completion scripts → completions/
	mkdir -p completions
	$(BINARY_CLI) completion bash > completions/goini.bash
	$(BINARY_CLI) completion zsh  > completions/goini.zsh
	$(BINARY_CLI) completion fish > completions/goini.fish

# ── Clean ─────────────────────────────────────────────────────────────────────

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin/ coverage.out
