# go-initializer — Implementation Task List

## Goal
Deliver go-initializer as both a web UI and a standalone `goini` CLI binary. The generation engine accepts a project specification (type, Go version, framework, addons, docker support, module name, project name, description) and produces a complete, immediately runnable Go project. The CLI binary is distributable via `go install` and Homebrew with no server or browser required. Tasks are ordered **low → high complexity** within each phase so implementation can proceed step by step without blockers.

---

## Motivation

The web UI is a good discovery surface, but developers live in the terminal. The core value of go-initializer — deterministic, reproducible project scaffolding — maps directly onto CLI workflows:

- Scriptable in CI/CD pipelines
- Works offline after installation
- Integrates with shell aliases and dotfiles
- Single binary, no runtime dependencies
- Installable in seconds via `go install` or `brew install`

---

## Current State Overview

### Backend
| File | Status |
|---|---|
| Server, router, config | Complete |
| `GET /api/meta` | Complete |
| `POST /api/generate` (routing) | Complete — registry-based, no switch statement |
| `registry.go` — `Generator` / `AddonGenerator` interfaces + registries | Complete |
| `gen_utils.go` — all helpers (`addToZip`, `GenerateGoModV2`, `GenerateMainContent`, etc.) | Complete — framework-aware (7 frameworks) |
| `gen_simple_project.go` — `SimpleProjectGenerator` | Complete |
| `gen_microservice.go` — `MicroserviceGenerator` | Complete |
| `gen_cli_app.go` — `CLIAppGenerator` | Complete |
| `gen_api_server.go` — `APIServerGenerator` | **Missing** |
| `router.go` — CORS | Fixed |
| `handler.go` | Complete |
| Input validation | Complete |

### Frontend
| File | Status |
|---|---|
| UI layout, form, theming | Complete |
| `useGeneratorForm.ts` hook | Complete |
| `service.ts` — `getMetaData` + `generateProject` | Complete |
| `GeneratorForm.tsx` — addon options | Complete — driven by API metadata |
| Error handling | Complete — inline error banner |
| Post-generate feedback | Complete — dismissible success banner |

### Known Infrastructure Issues
| Issue | Details |
|---|---|
| Backend port mismatch | `docker-compose.yml` maps `8181:8181` but server listens on `:8182` and Dockerfile exposes `8182` |
| `REACT_APP_API_URL` silently ignored in Docker | React bakes env vars at build time; `docker-compose.yml` sets the var at runtime, which has no effect — must be passed as a Docker `--build-arg` |
| `/api` prefix inconsistency | `frontend/.env` uses `http://localhost:8182/api` but `service.ts` fallback uses `http://localhost:8182` (no `/api` prefix) — calls would target non-existent routes |
| `config.go` dead code | `loadConfig` is defined but never called; server hardcodes port and timeouts |

---

## Phase 1 — Bug Fixes & Static Scaffolding Additions
> Low complexity. Fix broken code first, then add files that require no conditional logic.

- [x] **T1 — [Frontend] Fix API endpoint bug** — `service.ts` `generateProject` now correctly calls `/api/generate`
- [x] **T2 — [Backend] Fix `GenerateDatabaseAddon`** — replaced undefined `mysql`/`dsn` references; gorm uses `os.Getenv("DATABASE_URL")` + `postgres.Open`; ent uses env var DSN with commented-out scaffold
- [x] **T3 — [Backend] Add input validation** — add `validate` struct tags to `CreateProjectRequest` fields (`required`, `oneof` for enums); call the validator in `GenerateHandler` before processing; return `400` with field-level error messages on failure
- [x] **T4 — [Backend] Add `.gitignore` to all generated zips** — a `GenerateGitignore() []byte` helper that returns a standard Go `.gitignore` (binaries, `vendor/`, `.env`, IDE dirs); call it in every generator
- [x] **T5 — [Backend] Add `Makefile` to all generated zips** — a `GenerateMakefile(name string) []byte` helper with targets: `build`, `run`, `tidy` (`go mod tidy`), `test`; call it in every generator
- [x] **T6 — [Frontend] Add `logrus` to "other" addons display** — add `{ value: 'logrus', label: 'Logrus' }` to the `other` category in `addonOptions` inside `GeneratorForm.tsx`

---

## Phase 2 — Framework-Aware Code Generation (Core Engine)
> Medium complexity. The shared generation utilities must understand framework context before any specific generator can be completed.

- [x] **T7 — [Backend] Make `GenerateMainContent` framework-aware** — add a `framework string` parameter; return different `main.go` bodies per framework:
  - **golly** — existing `l3.Get()` + logger pattern
  - **gin** — `gin.Default()`, sample route, `router.Run(":8080")`
  - **echo** — `echo.New()`, sample route, `e.Start(":8080")`
  - **fiber** — `fiber.New()`, sample route, `app.Listen(":8080")`
  - **chi** — `chi.NewRouter()`, sample route, `http.ListenAndServe(":8080", r)`
  - **cobra** — `rootCmd.Execute()` entrypoint
  - **gokit** — minimal transport/endpoint/service entrypoint
- [x] **T8 — [Backend] Add `GenerateHandlerContent(framework string) []byte` helper** — produces a framework-specific `handler.go` with a sample health-check or hello handler; used by microservice and api-server generators
- [x] **T9 — [Backend] Add `GenerateRouterContent(framework string) []byte` helper** — produces a framework-specific `router.go` that wires routes to handlers; used by microservice and api-server generators
- [x] **T10 — [Backend] Add `GenerateServiceContent(name string) []byte` helper** — produces a simple `service.go` (package `internal`) with a stub interface and implementation; used by all generators
- [x] **T11 — [Backend] Add `GenerateLoggingAddon(addons []string) ([]byte, error)` helper** — produces an `internal/logger/logger.go` file for `zap` or `logrus` with an initialised logger instance; fixes the currently ignored "other" addons

---

## Phase 3 — Complete `simple-project` Generator
> Medium complexity. Builds directly on Phase 2 helpers; simplest project type.

- [x] **T12 — [Backend] Update `gen_simple_project.go` to use framework-aware `GenerateMainContent`** — pass `request.Framework` to `GenerateMainContent`; the only supported framework for simple-project is `golly` today but the call should be generic
- [x] **T13 — [Backend] Add logging addon support to `gen_simple_project.go`** — handle `addonType == "other"` in the addon loop; call `GenerateLoggingAddon` and write the result to `internal/logger/logger.go` in the zip

---

## Phase 4 — Complete `microservice` Generator
> Medium-high complexity. Multiple files, framework-specific structure, full addon support.

Project layout to generate:
```
<name>/
├── cmd/<name>/main.go        # framework-specific entrypoint
├── internal/
│   ├── handler.go            # framework-specific handler
│   ├── router.go             # framework-specific router
│   ├── service.go            # stub service interface + impl
│   ├── cache/cache.go        # if cache addon selected
│   ├── database/database.go  # if database addon selected
│   └── logger/logger.go      # if zap/logrus addon selected
├── go.mod
├── .gitignore
├── Makefile
├── README.md
└── Dockerfile                # if dockerSupport is true
```

- [x] **T14 — [Backend] Complete `GenerateMicroservice`** — the generator structure exists; wire Phase 2 helpers to make `cmd/main.go` framework-aware; add `internal/router.go` from `GenerateRouterContent`; replace the current hardcoded `internal/handler/handler.go` with `GenerateHandlerContent`; replace the current `internal/service/service.go` stub with `GenerateServiceContent`; handle "other" logging addons; emit `.gitignore` and `Makefile`

---

## Phase 5 — `api-server` Generator
> Medium-high complexity. Nearly identical layout to microservice but semantically a standalone HTTP API.

Project layout to generate:
```
<name>/
├── cmd/<name>/main.go        # framework-specific entrypoint
├── internal/
│   ├── handler.go
│   ├── router.go
│   ├── service.go
│   ├── cache/cache.go        # optional
│   ├── database/database.go  # optional
│   └── logger/logger.go      # optional
├── go.mod
├── .gitignore
├── Makefile
├── README.md
└── Dockerfile                # optional
```

Supported frameworks: `gin`, `echo`, `fiber`, `chi`, `golly`

- [ ] **T15 — [Backend] Create `gen_api_server.go`** — produce the full layout above using Phase 2 helpers; support all five frameworks and all addon types
- [ ] **T16 — [Backend] Register `api-server` in `generatorRegistry`** — add `"api-server": &APIServerGenerator{}` to `registry.go`; no changes to `handler.go` required

---

## Phase 6 — `cli-app` Generator
> High complexity. Different project layout; framework switch between cobra, urfave/cli, kingpin.

Project layout to generate:
```
<name>/
├── main.go                   # calls cmd.Execute()
├── cmd/
│   ├── root.go               # root command definition
│   └── <name>.go             # sample sub-command
├── internal/
│   └── logger/logger.go      # optional
├── go.mod
├── .gitignore
├── Makefile
└── README.md
```

Supported frameworks: `cobra`, `urfave`, `kingpin`

- [x] **T17 — [Backend] Add `GenerateRootCmd(framework, name string) []byte` helper** — returns `cmd/root.go` tailored to cobra / urfave / kingpin
- [x] **T18 — [Backend] Add `GenerateSubCmd(framework, name string) []byte` helper** — returns `cmd/<name>.go` with a sample sub-command for the chosen framework
- [x] **T19 — [Backend] Create `gen_cli_app.go`** — produce the full layout above using T17/T18 helpers and Phase 2 helpers; support all three CLI frameworks and logging addons
- [x] **T20 — [Backend] Register `cli-app` in `generatorRegistry`** — add `"cli-app": &CLIAppGenerator{}` to `registry.go`; no changes to `handler.go` required

---

## Phase 7 — Frontend Dynamic Addon UI & Polish
> Low-medium complexity. Frontend-only; depends on backend being stable.

- [x] **T21 — [Frontend] Drive addon UI from API metadata** — in `useGeneratorForm.ts`, derive addon options from the `supportedAddons` map returned by `getMetaData()`; expose them from the hook; remove the hardcoded `addonOptions` map in `GeneratorForm.tsx`
- [x] **T22 — [Frontend] Inline error display** — replace `alert(error.message)` in `useGeneratorForm.ts` with a state variable `generateError`; render an inline error banner in `GeneratorForm.tsx` below the generate button
- [x] **T23 — [Frontend] Post-generate success banner** — on successful blob download, set a `generateSuccess` state; render a dismissible success message in `GeneratorForm.tsx`

---

## Phase 8 — Infrastructure & DevX
> Low complexity. Fix pre-existing infra bugs before the repo restructure begins. No logic changes.

- [ ] **T24 — [Backend] Resolve dead config code** — `config.go` and `loadConfig` are never called; either wire `config.yaml` + call `loadConfig` in `Start()` to drive host/port/timeouts, or delete the dead code
- [ ] **T25 — [Infra] Fix `docker-compose.yml` port mismatch** — backend server listens on `:8182` and `Dockerfile` exposes `8182`, but `docker-compose.yml` maps `8181:8181`; align all to `8182:8182`
- [ ] **T25b — [Infra] Fix `REACT_APP_API_URL` not being applied in Docker** — React bakes env vars into the JS bundle at build time, not runtime; the `environment:` block in `docker-compose.yml` is silently ignored; pass `API_BASE_URL` as a Docker `--build-arg` in the frontend service instead
- [ ] **T25c — [Frontend] Align `/api` prefix across all config** — `frontend/.env` sets `http://localhost:8182/api` but `service.ts` fallback is `http://localhost:8182` (no `/api` prefix); standardise both to `http://localhost:8182/api` so all calls correctly target `/api/generate` and `/api/meta`

---

## Phase 9 — Repo Restructure (Shared Library Extraction)
> Medium complexity. Prerequisite for the CLI binary. Moves generator code out of `package main` into an importable internal package so both the HTTP server and the CLI share the same engine. The frontend is unaffected — it only communicates over HTTP.

### Problem

All generator code currently lives in `package main` inside `backend/`. This makes it impossible to import from a second binary (`cmd/goini/`) without duplicating the entire codebase.

### Solution

Extract the generation engine into `internal/generator/` — an importable package shared by both the HTTP server and the CLI binary.

### New Repo Layout

New repo layout after this phase:
```
go-initializer/
├── cmd/
│   ├── server/                    # HTTP server binary (replaces backend/)
│   │   └── main.go
│   └── goini/                     # CLI binary entry point
│       └── main.go
├── internal/
│   ├── generator/                 # shared generation engine (package generator)
│   │   ├── registry.go
│   │   ├── types.go
│   │   ├── gen_utils.go
│   │   ├── gen_simple_project.go
│   │   ├── gen_microservice.go
│   │   ├── gen_cli_app.go
│   │   ├── gen_api_server.go
│   │   └── response.go
│   └── server/                    # HTTP layer (thin wrapper over internal/generator)
│       ├── server.go
│       ├── router.go
│       └── handler.go
├── go.mod                         # single root go.mod (replaces backend/go.mod)
├── go.sum
├── docker-compose.yml
└── frontend/                      # unchanged
```

- [x] **T26 — [Repo] Create root `go.mod`** — move `backend/go.mod` to repo root; update module path to `github.com/neo7337/go-initializer`; add `github.com/spf13/cobra` and `github.com/charmbracelet/huh` as new direct dependencies
- [x] **T27 — [Repo] Extract `internal/generator` package** — copy all generator files from `backend/` into `internal/generator/`; change `package main` → `package generator`; export all types and functions needed externally (registries, maps, helpers, `CreateProjectRequest`)
- [x] **T28 — [Repo] Migrate HTTP server to `internal/server` + `cmd/server`** — move `server.go`, `router.go`, `handler.go` into `internal/server/`; update all imports to use `internal/generator`; create `cmd/server/main.go` as the thin entrypoint (`package main; func main() { server.Start() }`)
- [x] **T29 — [Repo] Delete `backend/` and update all infra references** — after verifying `cmd/server` builds cleanly, remove `backend/`; update the following hardcoded `./backend` paths:
  - `docker-compose.yml` — `build.context` → `.` with `build.dockerfile: cmd/server/Dockerfile`
  - `.github/workflows/backend.yml` — `paths` filter → `cmd/server/**` and `internal/**`; `working-directory` → repo root; build command → `go build ./cmd/server/...`
  - `.github/workflows/release.yml` — change detection filter and Docker build context from `backend/` → `cmd/server/`

---

## Phase 10 — `goini` CLI Binary
> High complexity. Core user-facing deliverable of the CLI-first distribution model.

The `goini new` command scaffolds a project interactively (step-by-step prompts) with optional flags to pre-fill any answer — making it both discoverable for new users and scriptable in CI. Any flag provided skips its corresponding prompt.

**Commands:**
```
goini new                       # interactive project scaffold
goini list types                # print all supported project types
goini list frameworks           # print frameworks for a type (--type flag)
goini list addons               # print all supported addons
goini version                   # print version and build info
goini completion bash|zsh|fish  # print shell completion script
```

**`goini new` interactive flow:**
```
$ goini new

  Project name       myapp
  Module path        github.com/acme/myapp
  Description        (optional — press Enter to skip)
  Go version       › 1.25.0   1.24.6   1.23.12
  Project type     › Microservice   Simple Project   CLI Application   API Server
  Framework        › Gin   Echo   Fiber   Chi   Golly   (filtered by project type)
  Cache addon        [ ] Redis   [ ] Memcached
  Database addon     [ ] GORM    [ ] Ent
  Other addons       [ ] Zap     [ ] Logrus   [ ] Cobra
  Docker support     (y/N)
  Output directory   ./myapp

  Generating...
  Project created at ./myapp/

  Next steps:
    cd myapp
    go mod tidy
    make run
```

**`goini new` flags (any flag skips its prompt):**

| Flag | Type | Description |
|---|---|---|
| `--name` | string | Project name |
| `--module` | string | Go module path (e.g. `github.com/acme/myapp`) |
| `--description` | string | Short project description |
| `--go-version` | string | Go version — one of the values from `SupportedGoVersionsMap` |
| `--type` | string | `microservice`, `simple-project`, `cli-app`, `api-server` |
| `--framework` | string | Framework name — must be valid for the selected project type |
| `--addon` | string (repeatable) | `category=value` format, e.g. `--addon cache=redis --addon database=gorm` |
| `--docker` | bool | Generate a multi-stage `Dockerfile` |
| `--output` | string | Output directory (default: `./<name>`) |

Any flag provided skips its corresponding interactive prompt — making the command fully scriptable in CI:

```bash
goini new \
  --name        myapp \
  --module      github.com/acme/myapp \
  --description "A sample microservice" \
  --go-version  1.25.0 \
  --type        microservice \
  --framework   gin \
  --addon       cache=redis \
  --addon       database=gorm \
  --addon       other=zap \
  --docker \
  --output      ./projects/
```

**Output:** the generator returns a `*bytes.Buffer` (zip); the CLI:

1. Receives the zip buffer from `generator.generatorRegistry[projectType].Generate(request)`
2. Creates the output directory if it does not exist
3. Extracts all zip entries directly into the output directory using `archive/zip`
4. Does **not** write a `.zip` file to disk — developers want a folder, not an archive
5. Prints a success message with the absolute output path and contextual next-step hints

**Next-step hints by project type:**

| Project type | Hint |
|---|---|
| `microservice` / `api-server` / `simple-project` | `cd <name> && go mod tidy && make run` |
| `cli-app` | `cd <name> && go mod tidy && make build` |

- [x] **T30 — [CLI] Create `cmd/goini/main.go` skeleton** — cobra root command with `--version` flag and subcommand registration; wire `goini version` and `goini completion`
- [x] **T31 — [CLI] Implement `goini list` subcommands** — `list types`, `list frameworks --type <t>`, `list addons`; read directly from `internal/generator` registry maps; output as a formatted table
- [x] **T32 — [CLI] Implement `goini new` flag parsing** — define all flags listed above; collect provided values into a partial `generator.CreateProjectRequest`
- [x] **T33 — [CLI] Implement interactive prompts for `goini new`** — use `github.com/charmbracelet/huh` to prompt for any field not supplied via flag; filter framework choices by selected project type using `SupportedFrameworksMap`
- [ ] **T34 — [CLI] Wire `goini new` to the generation engine** — build a complete `generator.CreateProjectRequest` from prompts + flags; call `generator.generatorRegistry[projectType].Generate(request)`; extract the returned zip buffer into the output directory using `archive/zip`
- [ ] **T35 — [CLI] Add post-generation output** — print success message with absolute output path and contextual next-step hints (`make run` for server types, `make build` for CLI types)

---

## Phase 11 — Distribution & Release
> Low-medium complexity. Packaging and delivery pipeline.

### `go install` (Primary)

Once the repo is restructured with a root `go.mod`, the CLI installs with the standard Go toolchain — no additional tooling required:

```bash
go install github.com/neo7337/go-initializer/cmd/goini@latest
```

This is the recommended installation path for Go developers.

### Homebrew (Secondary)

A dedicated tap repository (`github.com/neo7337/homebrew-goini`) enables installation without the Go toolchain:

```bash
brew tap neo7337/goini
brew install goini
```

Shell completions are installed automatically by the formula.

### Cross-Compilation with goreleaser

`goreleaser` handles cross-compilation and release asset management. Add `goreleaser.yaml` at the repo root:

```yaml
# goreleaser.yaml (outline)
builds:
  - id: goini
    main: ./cmd/goini
    binary: goini
    goos: [linux, darwin]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w
      - -X main.version={{.Version}}

archives:
  - id: goini
    builds: [goini]
    name_template: "goini_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    files:
      - completions/*

brews:
  - name: goini
    repository:
      owner: neo7337
      name: homebrew-goini
    install: |
      bin.install "goini"
      bash_completion.install "completions/goini.bash" => "goini"
      zsh_completion.install "completions/goini.zsh" => "_goini"
      fish_completion.install "completions/goini.fish"
```

### Shell Completions

Cobra generates completion scripts automatically. A `Makefile` target generates them at release time:

```bash
goini completion bash > completions/goini.bash
goini completion zsh  > completions/goini.zsh
goini completion fish > completions/goini.fish
```

Manual installation for users not using Homebrew:

```bash
# Bash
goini completion bash > /usr/local/etc/bash_completion.d/goini

# Zsh
goini completion zsh > ~/.zsh/completions/_goini

# Fish
goini completion fish > ~/.config/fish/completions/goini.fish
```

### GitHub Actions Release Pipeline

Extend `.github/workflows/release.yml`:

```
On: GitHub Release published
  1. Run goreleaser release
     → Cross-compiles goini for all targets
     → Attaches binaries to the GitHub Release
     → Updates Homebrew formula with new SHA256 + version
  2. Build and push cmd/server Docker image to Docker Hub
     (separate job, existing behaviour preserved)
```

- [ ] **T36 — [Infra] Add `goreleaser.yaml`** — cross-compile `goini` for `darwin/amd64`, `darwin/arm64`, `linux/amd64`, `linux/arm64`; attach binaries to GitHub Releases; include shell completion scripts in archives; configure `brews` block to auto-update the Homebrew formula on each release
- [ ] **T37 — [Infra] Update `.github/workflows/release.yml`** — run `goreleaser release` on GitHub Release event; preserve the existing Docker Hub publish job for the server image (`cmd/server`)
- [ ] **T38 — [Infra] Create Homebrew tap** — create `github.com/neo7337/homebrew-goini` repository with `Formula/goini.rb`; formula installs the binary and all three shell completion scripts (bash, zsh, fish)
- [ ] **T39 — [Docs] Update README** — add `go install` one-liner, `brew install` command, `goini new` quick-start example, and a link to the web UI for users who prefer a graphical interface

---

## Execution Order Summary

| Phase | Tasks | Complexity | Dependency | Status |
|---|---|---|---|---|
| 1 — Bug Fixes & Static Files | T1–T6 | Low | None | All done |
| 2 — Framework-Aware Engine | T7–T11 | Medium | Phase 1 | All done |
| 3 — Complete `simple-project` | T12–T13 | Medium | Phase 2 | All done |
| 4 — Complete `microservice` | T14 | Medium-High | Phase 2 | All done |
| 5 — `api-server` Generator | T15–T16 | Medium-High | Phase 2 | Open |
| 6 — `cli-app` Generator | T17–T20 | High | Phase 2 | All done |
| 7 — Frontend Dynamic UI | T21–T23 | Low-Medium | Phase 1 | All done |
| 8 — Infrastructure & DevX | T24–T25c | Low | Any | Open |
| 9 — Repo Restructure | T26–T29 | Medium | Phase 5 + Phase 8 | Open |
| 10 — `goini` CLI Binary | T30–T35 | High | Phase 9 | Open |
| 11 — Distribution & Release | T36–T39 | Low-Medium | Phase 10 | Open |

---

## Dependencies to Add

| Package | Version | Purpose |
|---|---|---|
| `github.com/spf13/cobra` | `v1.10.2` | CLI command framework for `goini` |
| `github.com/charmbracelet/huh` | `latest` | Terminal form / interactive prompts |

Both packages are pure Go with no CGO requirements — the binary remains fully static and cross-compilable.

---

## Implementation Order

| Step | Phase | Unblocks |
|---|---|---|
| 1. Complete `api-server` generator (T15, T16) | Phase 5 | Phase 9 (clean extraction) |
| 2. Root `go.mod` + extract `internal/generator` (T26, T27) | Phase 9 | Everything CLI |
| 3. Migrate HTTP server to `cmd/server` (T28) | Phase 9 | T29 |
| 4. Delete `backend/`, update CI (T29) | Phase 9 | Phase 10 |
| 5. CLI skeleton + `list` commands (T30, T31) | Phase 10 | T32–T35 |
| 6. `goini new` flags + prompts + engine wiring (T32–T35) | Phase 10 | Phase 11 |
| 7. goreleaser + release workflow + Homebrew (T36–T38) | Phase 11 | T39 |
| 8. README update (T39) | Phase 11 | — |
