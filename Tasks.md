# go-initializer — Task Backlog
Derived from gap analysis against today's Go ecosystem and developer tooling market.
Tasks are grouped by priority and category. Each task includes context on why it matters.
---
## Priority 1 — Critical / Broken
### T-001: Implement `api-server` generator
**File:** `internal/generator/` (new file `gen_api_server.go`)
**Why:** `api-server` is registered in `SupportedProjectTypesMap` and `SupportedFrameworksMap` but has no entry in `GeneratorRegistry`. Selecting it returns a runtime error. This is a broken promise in the README and the web UI.
**Acceptance criteria:**
- [ ] New `APIServerGenerator` struct implementing the `Generator` interface
- [ ] Registered in `GeneratorRegistry` under `"api-server"`
- [ ] Supports all five declared frameworks: `gin`, `echo`, `fiber`, `chi`, `golly`
- [ ] Generated layout mirrors the microservice layout (handler, router, service layers)
- [ ] All addons (cache, database, logging, AI, vectorstore) work with the new type
- [ ] Unit tests added in `generators_test.go`
---
## Priority 2 — High Impact / Competitive Gap
### T-002: Add gRPC / protobuf project scaffold
**Why:** gRPC is the dominant microservice-to-microservice transport in Go. Absence is a major competitive gap vs. tools like `buf`. No other Go scaffolder covers this well.
**Acceptance criteria:**
- [ ] New framework option `grpc` under `microservice` and `api-server` project types
- [ ] Generated layout includes a sample `.proto` file, a `buf.yaml`, and `buf.gen.yaml`
- [ ] Generated `server.go` registers a gRPC server with a health check service
- [ ] `go.mod` wired with `google.golang.org/grpc` and `google.golang.org/protobuf`
- [ ] Makefile includes a `proto` target that runs `buf generate`
- [ ] README explains the `buf generate` step
### T-003: Add Connect (connectrpc.com) framework support
**Why:** Connect is the modern HTTP/1.1 + HTTP/2 gRPC-compatible protocol gaining rapid adoption in 2024-25. Used at Buf, PlanetScale, and others.
**Acceptance criteria:**
- [ ] New framework option `connect` under `microservice` and `api-server` project types
- [ ] Generated layout includes a `.proto` file and a Connect handler
- [ ] `go.mod` wired with `connectrpc.com/connect`
- [ ] Sample health endpoint using the Connect protocol
### T-004: Add GraphQL (`gqlgen`) scaffold
**Why:** GraphQL is a top request in any API-oriented Go scaffolder. `gqlgen` is the dominant Go GraphQL library.
**Acceptance criteria:**
- [ ] New framework option `gqlgen` under `api-server` project type
- [ ] Generated layout includes a `schema.graphqls`, `gqlgen.yml`, and resolver stubs
- [ ] `go.mod` wired with `github.com/99designs/gqlgen`
- [ ] Generated `main.go` starts the GraphQL server with a playground handler
### T-005: Add Anthropic Claude to AI providers
**Why:** Claude is one of the two most widely used LLM APIs as of 2025. Absence is a notable gap for an AI-first scaffolder.
**Acceptance criteria:**
- [ ] New framework option `anthropic` under `ai-agent` project type
- [ ] New addon value `anthropic` under the `ai` addon category
- [ ] `DependencyMap` entry for `github.com/anthropics/anthropic-sdk-go`
- [ ] `GenerateLLMClient` case for `anthropic` (llm/client.go)
- [ ] `GenerateAgentContent` case for `anthropic` (agent/agent.go — tool-use loop)
- [ ] `GenerateToolsContent` case for `anthropic` (tools/tools.go — tool definitions)
- [ ] `GenerateAIAddonContent` case for `anthropic` (internal/ai/client.go)
### T-006: Add OpenTelemetry observability addon
**Why:** OpenTelemetry is the de-facto tracing and metrics standard in 2025. Every serious microservice template includes it. Currently no observability option exists.
**Acceptance criteria:**
- [ ] New addon category `observability` with value `otel`
- [ ] Registered in `SupportedAddonsMap` and `addonRegistry`
- [ ] Generated `internal/telemetry/telemetry.go` initialises an OTLP trace exporter and a Prometheus metrics exporter
- [ ] `DependencyMap` entries for `go.opentelemetry.io/otel` and related packages
- [ ] `go.mod` includes all otel dependencies when selected
- [ ] Works across microservice, api-server, and ai-agent project types
### T-007: Add Prometheus metrics addon
**Why:** `prometheus/client_golang` is the default metrics library for Go services and the most commonly added dependency after a web framework.
**Acceptance criteria:**
- [ ] New addon value `prometheus` under a `observability` (or `other`) category
- [ ] Generated `internal/metrics/metrics.go` registers a default metrics registry and a `/metrics` handler
- [ ] `DependencyMap` entry for `github.com/prometheus/client_golang`
- [ ] `/metrics` route wired into the framework-specific router when selected
### T-008: Add `slog` (stdlib) as a logging option
**Why:** `slog` was added in Go 1.21 and is now the idiomatic structured logging choice. It should be the first option listed, before `zap` and `logrus`.
**Acceptance criteria:**
- [ ] New addon value `slog` under the `other` addon category
- [ ] `GenerateLoggingAddon` case for `slog` (no external dependency — stdlib only)
- [ ] Generated `internal/logger/logger.go` sets up a JSON `slog.Logger` backed by `os.Stdout`
- [ ] `SupportedAddonsMap` updated
### T-009: Add GitHub Actions CI workflow to generated projects
**Why:** CI is expected in every new repo. A generated `.github/workflows/ci.yml` that runs `make test` and `make build` removes a manual step that most developers would otherwise copy-paste.
**Acceptance criteria:**
- [ ] Generated `.github/workflows/ci.yml` for all project types
- [ ] Workflow runs on `push` and `pull_request` to `main`
- [ ] Steps: checkout, setup-go (using the selected Go version), `make tidy`, `make build`, `make test`
- [ ] Optional: add `golangci-lint` step when lint config is also generated (see T-015)
---
## Priority 3 — Medium Impact / Completeness
### T-010: Add `sqlc` database addon
**Why:** `sqlc` (type-safe SQL codegen) is very popular in the Go community in 2024-25 and is the preferred alternative to ORMs for many developers.
**Acceptance criteria:**
- [ ] New addon value `sqlc` under the `database` category
- [ ] Generated `sqlc.yaml` configuration file
- [ ] Sample `query.sql` and `schema.sql` in a `db/` directory
- [ ] Generated `internal/database/database.go` using the `database/sql` stdlib with `pgx` driver
- [ ] `DependencyMap` entry for `github.com/sqlc-dev/sqlc` and `github.com/jackc/pgx/v5`
- [ ] Makefile `generate` target runs `sqlc generate`
### T-011: Add MongoDB database addon
**Why:** MongoDB is very common in Go services. Its official driver is mature and widely used.
**Acceptance criteria:**
- [ ] New addon value `mongodb` under the `database` category
- [ ] Generated `internal/database/database.go` with `mongo.Connect` and a ping health check
- [ ] `DependencyMap` entry for `go.mongodb.org/mongo-driver/v2`
- [ ] Connection string loaded from `MONGODB_URI` environment variable
### T-012: Expose `zerolog` as a logging addon
**Why:** `zerolog` is already present in `DependencyMap` (`registry.go:129`) but is not in `SupportedAddonsMap` — it is unreachable dead code. It is one of the most popular Go logging libraries.
**Acceptance criteria:**
- [ ] Add `zerolog` to `SupportedAddonsMap["other"]`
- [ ] `GenerateLoggingAddon` case for `zerolog`
- [ ] Generated `internal/logger/logger.go` sets up a `zerolog.Logger` with a console writer
### T-013: Expose `viper` as a config addon
**Why:** `viper` is already in `DependencyMap` (`registry.go:130`) but not in `SupportedAddonsMap` — unreachable dead code. Config management is a common need alongside any framework.
**Acceptance criteria:**
- [ ] New addon category `config` with value `viper`
- [ ] Registered in `SupportedAddonsMap` and `addonRegistry`
- [ ] Generated `internal/config/config.go` that loads from env vars and an optional `config.yaml`
- [ ] Sample `config.yaml` included in the generated project
- [ ] `DependencyMap` entry confirmed / used
### T-014: Generate `docker-compose.yml` for the project's dev dependencies
**Why:** The tool itself has a `docker-compose.yml` but generated projects get none. A user who selects Redis + Postgres addons should get a compose file that starts those services locally.
**Acceptance criteria:**
- [ ] `docker-compose.yml` generated when any infrastructure addon (cache, database) is selected
- [ ] Services reflect the selected addons (e.g. `redis:7-alpine` for Redis, `postgres:17-alpine` for GORM/ent/pgvector)
- [ ] Environment variables in the compose file match those expected by the generated addon code (`DATABASE_URL`, `REDIS_ADDR`, etc.)
- [ ] Generated only when at least one infrastructure addon is selected
### T-015: Add `golangci-lint` config to generated projects
**Why:** `.golangci.yml` is standard in production Go repos. Linting setup is expected by teams adopting a new project scaffold.
**Acceptance criteria:**
- [ ] Generated `.golangci.yml` with a sensible default ruleset (`errcheck`, `govet`, `staticcheck`, `unused`, `gosimple`)
- [ ] Makefile `lint` target: `golangci-lint run ./...`
- [ ] `.github/workflows/ci.yml` (T-009) includes a lint step
### T-016: Add `air` live-reload config to generated projects
**Why:** `air` is the standard live-reload tool for Go web development. An `.air.toml` prevents developers from having to look up the config on day one.
**Acceptance criteria:**
- [ ] Generated `.air.toml` for microservice, api-server, and ai-agent project types
- [ ] `build.cmd` points to the correct main package path
- [ ] `build.bin` points to the binary output path matching the Makefile
- [ ] Makefile `dev` target: `air`
### T-017: Generate `.env.example` for projects requiring secrets
**Why:** Projects generating LLM clients, database connections, or cache connections require environment variables at startup. Without `.env.example`, users have no reference for what to set.
**Acceptance criteria:**
- [ ] `.env.example` generated for any project that requires env vars (AI, database, cache addons)
- [ ] File lists all expected variables with placeholder values and a one-line comment per variable
- [ ] `.gitignore` already excludes `.env` (confirmed) — `.env.example` must NOT be excluded
### T-018: Fix multi-addon-per-category silently dropping extras
**Why:** When multiple addons in the same category are selected, only the first is generated (`gen_ai_addon.go:17`, `gen_ai_agent.go:121`). The second selection is silently ignored with no warning to the user.
**Acceptance criteria:**
- [ ] Either generate all selected addons in a category (preferred), or
- [ ] Return a validation error / warning when more than one addon in a mutually-exclusive category is selected
- [ ] Web UI and CLI clearly communicate single-select vs multi-select constraints per category
- [ ] Existing tests updated; new test covering multi-addon selection
### T-019: Add generated handler and service unit tests
**Why:** Generated projects have a `make test` target that runs against empty packages. No test files are generated, making the test target a no-op.
**Acceptance criteria:**
- [ ] Generated `internal/handler/handler_test.go` with a basic HTTP handler test using `net/http/httptest`
- [ ] Generated `internal/service/service_test.go` with a basic interface-level test
- [ ] Tests pass immediately after `go mod tidy` without modification
- [ ] `testify` added to `go.mod` when test files are generated
### T-020: Add Kubernetes manifests to generated projects
**Why:** The vast majority of production Go services target Kubernetes. Basic `Deployment`, `Service`, and `ConfigMap` manifests remove significant boilerplate.
**Acceptance criteria:**
- [ ] New optional flag `--kubernetes` / `kubernetesSupport bool` in `CreateProjectRequest`
- [ ] Generated `k8s/deployment.yaml`, `k8s/service.yaml`, `k8s/configmap.yaml`
- [ ] Manifests reference the correct container image name matching the project name
- [ ] `ConfigMap` keys match the env vars expected by the generated code
- [ ] Web UI exposes the option alongside the existing Docker toggle
---
## Priority 4 — Polish / Developer Experience
### T-021: Add `sqlx` as an exposed database addon
**Why:** `sqlx` is already in `DependencyMap` (`registry.go:126-127`) but not in `SupportedAddonsMap` — unreachable dead code. It is a popular lightweight alternative to GORM.
**Acceptance criteria:**
- [ ] Add `sqlx` to `SupportedAddonsMap["database"]`
- [ ] `GenerateDatabaseAddon` case for `sqlx`
- [ ] Generated `internal/database/database.go` with `sqlx.Connect` and a ping health check
- [ ] Connection string loaded from `DATABASE_URL`
### T-022: Support multiple Postgres drivers in GORM addon
**Why:** GORM's generated code is hardcoded to `gorm.io/driver/postgres`. Users targeting MySQL, SQLite, or SQL Server get unusable code.
**Acceptance criteria:**
- [ ] GORM addon accepts a sub-option for driver: `postgres` (default), `mysql`, `sqlite`
- [ ] Generated `database.go` imports and uses the correct driver package
- [ ] `DependencyMap` includes entries for `gorm.io/driver/mysql` and `gorm.io/driver/sqlite`
### T-023: Add `godotenv` loading to generated projects using secrets
**Why:** Generated projects that require `OPENAI_API_KEY`, `DATABASE_URL`, etc. will panic on startup in a local dev environment without those vars set. Loading a `.env` file prevents a confusing first-run experience.
**Acceptance criteria:**
- [ ] When any addon requiring env vars is selected, `github.com/joho/godotenv` is added to `go.mod`
- [ ] Generated `main.go` (or the appropriate init path) calls `godotenv.Load()` at startup
- [ ] Load is best-effort (non-fatal if `.env` is absent) — appropriate for production use
### T-024: Improve Dockerfile to use correct binary path per project type
**Why:** The generated `Dockerfile` always runs `go build -o main .` which does not match the project layout for microservices (where `main` is at `cmd/<name>/`) or CLI apps.
**Acceptance criteria:**
- [ ] `GenerateDockerfile` uses `request.ProjectType` and `request.Name` to set the correct build path
- [ ] Microservice / api-server: `go build -o bin/<name> ./cmd/<name>`
- [ ] CLI app: `go build -o bin/<name> .`
- [ ] AI agent: `go build -o bin/<name> .`
- [ ] CMD in final stage uses the correct binary path
### T-025: Add WebSocket / SSE scaffold for real-time transports
**Why:** Streaming LLM responses and real-time event feeds are common. No scaffold for WebSocket or SSE exists, even in the AI agent project type where streaming is most relevant.
**Acceptance criteria:**
- [ ] New addon value `websocket` under a `transport` or `other` category
- [ ] Generated `internal/ws/handler.go` with a basic echo WebSocket handler using `gorilla/websocket` or `nhooyr.io/websocket`
- [ ] Alternatively, an SSE handler option for AI agent streaming responses
- [ ] Wired into the framework-specific router when selected
### T-026: Add `testcontainers-go` integration test scaffold
**Why:** Integration tests against real databases/caches are the standard. `testcontainers-go` is the dominant library for this pattern in Go.
**Acceptance criteria:**
- [ ] When a database or cache addon is selected, generate an optional `internal/database/integration_test.go`
- [ ] Uses `testcontainers-go` to spin up the correct container (postgres, redis, etc.)
- [ ] Test verifies the connection can be established
- [ ] `DependencyMap` entry for `github.com/testcontainers/testcontainers-go`
### T-027: Surface Go version in generated README next-steps
**Why:** Generated `README.md` files contain only `# <name>\n\n<description>`. They don't mention the selected Go version, framework, or addons, making them less useful as onboarding docs.
**Acceptance criteria:**
- [ ] Generated README includes: Go version, project type, framework, selected addons
- [ ] Next-steps section is framework-aware (e.g. gRPC README mentions `make proto`, AI agent README mentions setting API keys)
- [ ] README is generated consistently across all project types
### T-028: Add `buf` toolchain support for gRPC projects (dependency on T-002)
**Why:** `buf` is the modern replacement for raw `protoc` calls. Any gRPC scaffold should use `buf` rather than manual `protoc` invocations.
**Acceptance criteria:**
- [ ] Generated `buf.yaml` (module config) and `buf.gen.yaml` (code generation config)
- [ ] `buf.gen.yaml` configured to generate Go and gRPC stubs into a `gen/` directory
- [ ] Makefile `proto` target uses `buf generate` (not raw `protoc`)
- [ ] `.gitignore` excludes the `gen/` directory (since it is derived from `.proto` files)
---
## Housekeeping
### T-029: Remove or expose dead code in `DependencyMap`
**Why:** `DependencyMap` contains entries for `zerolog`, `viper`, `sqlx`, `httptest`, and `testify` that are never reachable through any supported addon or project type. This is confusing for contributors.
**Acceptance criteria:**
- [ ] Each entry in `DependencyMap` is either reachable via a registered addon/framework, or removed
- [ ] `zerolog` → resolved by T-012
- [ ] `viper` → resolved by T-013
- [ ] `sqlx` → resolved by T-021
- [ ] `testify` / `httptest` → resolved by T-019
### T-030: Add `go.sum` guidance / stub to generated projects
**Why:** Generated `go.mod` is never accompanied by a `go.sum`. Users must run `go mod tidy` immediately after generation. If the network is unavailable or a dependency version has been retracted, the first run fails confusingly.
**Acceptance criteria:**
- [ ] Generated `README.md` prominently lists `go mod tidy` as the mandatory first step (above `make run`)
- [ ] Alternatively, investigate generating a minimal `go.sum` at generation time by resolving the dependency graph server-side
- [ ] The web UI "what to do next" panel explicitly calls out `go mod tidy`
---
## Task Summary
| ID | Title | Priority |
|---|---|---|
| T-001 | Implement `api-server` generator | Critical |
| T-002 | Add gRPC / protobuf scaffold | High |
| T-003 | Add Connect framework support | High |
| T-004 | Add GraphQL (`gqlgen`) scaffold | High |
| T-005 | Add Anthropic Claude AI provider | High |
| T-006 | Add OpenTelemetry observability addon | High |
| T-007 | Add Prometheus metrics addon | High |
| T-008 | Add `slog` logging option | High |
| T-009 | Generate GitHub Actions CI workflow | High |
| T-010 | Add `sqlc` database addon | Medium |
| T-011 | Add MongoDB database addon | Medium |
| T-012 | Expose `zerolog` logging addon | Medium |
| T-013 | Expose `viper` config addon | Medium |
| T-014 | Generate `docker-compose.yml` for dev deps | Medium |
| T-015 | Add `golangci-lint` config | Medium |
| T-016 | Add `air` live-reload config | Medium |
| T-017 | Generate `.env.example` | Medium |
| T-018 | Fix multi-addon silent drop bug | Medium |
| T-019 | Generate handler and service unit tests | Medium |
| T-020 | Add Kubernetes manifests option | Medium |
| T-021 | Expose `sqlx` database addon | Polish |
| T-022 | Support multiple GORM drivers | Polish |
| T-023 | Add `godotenv` loading for secret-requiring projects | Polish |
| T-024 | Fix Dockerfile binary path per project type | Polish |
| T-025 | Add WebSocket / SSE scaffold | Polish |
| T-026 | Add `testcontainers-go` integration test scaffold | Polish |
| T-027 | Improve generated README with project metadata | Polish |
| T-028 | Add `buf` toolchain support for gRPC | Polish |
| T-029 | Remove or expose dead code in `DependencyMap` | Housekeeping |
| T-030 | Add `go.sum` guidance to generated projects | Housekeeping |
TASKS.md has been created with 30 tasks across 4 priority tiers. Here's the structure:
Critical (1 task)
- T-001 — api-server generator is completely missing/broken at runtime
High Impact (8 tasks)
- T-002 through T-009 covering gRPC, Connect, GraphQL, Anthropic, OpenTelemetry, Prometheus, slog, and GitHub Actions CI
Medium Impact (11 tasks)
- T-010 through T-020 covering sqlc, MongoDB, dead-code logging addons (zerolog, viper), docker-compose for generated projects, golangci-lint, air, .env.example, the silent multi-addon bug, test generation, and Kubernetes manifests
Polish / Housekeeping (10 tasks)
- T-021 through T-030 covering sqlx, multi-driver GORM, godotenv, Dockerfile binary path fix, WebSocket/SSE, testcontainers, README improvements, buf toolchain, and dead-code cleanup
Each task has a clear rationale, references the specific file/line where the gap exists, and has acceptance criteria so a contributor knows exactly when it's done.