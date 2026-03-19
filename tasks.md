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
- [x] **T34 — [CLI] Wire `goini new` to the generation engine** — build a complete `generator.CreateProjectRequest` from prompts + flags; call `generator.generatorRegistry[projectType].Generate(request)`; extract the returned zip buffer into the output directory using `archive/zip`
- [x] **T35 — [CLI] Add post-generation output** — print success message with absolute output path and contextual next-step hints (`make run` for server types, `make build` for CLI types)

---

## Phase 11 — Frontend UI/UX Overhaul
> Medium-high complexity. Migrate off the unmaintained CRA build tool, establish a unified design system, and redesign every layer of the UI with a terminal/IDE-inspired dark-first aesthetic using the Geist font. All logic (hooks, services, types) is preserved — this phase is styling, structure, and build tooling only.

**Design direction:** Dark-first, terminal/IDE-inspired (Vercel/Railway aesthetic). Accent colour `#ffd700` (gold) formalized as a design token. Geist font via `@fontsource/geist`.

### Phase A — Foundation

- [x] **T-UX1 — [Frontend] Migrate CRA → Vite** — install `vite` and `@vitejs/plugin-react`; create `vite.config.ts`; move `public/index.html` to repo root and update the script tag; replace `react-scripts` with `vite`/`vitest` in `package.json`; rename all `REACT_APP_*` env vars to `VITE_*`; update `service.ts` to use `import.meta.env.VITE_API_URL`; update `frontend/Dockerfile` build command and nginx output path from `build/` to `dist/`

- [x] **T-UX2 — [Frontend] Design system token layer** — replace all fragmented CSS variables in `App.css` and hard-coded hex values throughout components with a unified token set: `--color-surface-*` (3 elevation levels), `--color-text-*` (primary / secondary / muted), `--color-border`, `--color-accent` (`#ffd700` formalized), `--color-success`, `--color-destructive`, `--radius-*`, `--shadow-*`; add Geist font via `@fontsource/geist` and apply it to `body`; set dark as the default appearance; remove unused `App-logo` keyframes and dead CSS

### Phase B — Layout & Navigation

- [x] **T-UX3 — [Frontend] App shell redesign** — replace all inline `React.CSSProperties` in `App.tsx` with Tailwind/CSS classes; new sticky header with backdrop-blur, Geist logo wordmark, theme toggle using SVG icons (not emoji), GitHub icon button with `aria-label`; slim footer (copyright + links only); centred max-width container for main content

### Phase C — Generator Form

- [x] **T-UX4 — [Frontend] Form visual redesign** — replace bare `<input>` elements with Radix UI `TextField`; replace radio buttons and checkboxes with styled pill/badge selectors for version and framework choices (larger hit targets, clear selected-state ring); render addon choices as tag-chip multi-select groups; add subtle numbered section labels (`01 Go Version`, `02 Project Type`, …); `GENERATE` button: full-width, prominent, uses `--color-accent` token, shows spinner during generation

- [x] **T-UX5 — [Frontend] Inline validation & feedback states** — redesign the error callout as a terminal-style inline banner with an icon; redesign the success callout with an animated checkmark and auto-dismiss after 5 s with a countdown progress bar; add inline field-level validation: red underline + message below each input (replaces border-colour-only feedback)

### Phase D — Explore / Docs Browser

- [x] **T-UX6 — [Frontend] Docs browser redesign** — replace the brittle positioned-div sidebar tree with a semantic `<nav>` element; add CSS `max-height` collapse/expand animation; active page: left accent-border + tinted background; replace emoji folder/file icons with inline SVG icons; add a breadcrumb bar at the top of the content pane; add syntax highlighting to fenced code blocks via `react-syntax-highlighter` + Prism; replace the "Loading…" text with a skeleton loader; add a copy-to-clipboard button on every code block

### Phase E — Polish & Quality

- [x] **T-UX7 — [Frontend] Responsive design** — add breakpoints for mobile (`< 768 px`) and tablet (`768–1024 px`): stack form sections vertically on mobile; collapse Explore sidebar into a hamburger/drawer on mobile; full-width inputs below `768 px`

- [x] **T-UX8 — [Frontend] Theme persistence + system preference** — persist theme choice to `localStorage`; honour `prefers-color-scheme` on first load when no saved preference exists; add `transition: background-color 200ms, color 200ms` for a smooth toggle

- [x] **T-UX9 — [Frontend] Accessibility** — `aria-label` on all icon-only buttons; `role="main"` on the content area; every `<input>` has an associated `<label>`; `focus-visible` outline using `--color-accent` token; full keyboard navigation through the form without a mouse

- [x] **T-UX10 — [Frontend] Micro-animations & polish** — fade-in on initial page load; card hover: subtle `box-shadow` lift using `--shadow-*` token; smooth transition when framework options change after a project-type selection; button loading spinner replaces button label text during generation

---

## Phase 12 — AI Agent Capabilities
> High complexity. Adds a new `ai-agent` project type and two new cross-cutting addon categories (`ai`, `vectorstore`) available to all project types. The `framework` field maps to LLM provider for `ai-agent`. Frontend picks up all changes automatically via `/api/meta`.

**Supported LLM providers (framework values for `ai-agent`):** `langchaingo`, `openai`, `gemini`, `ollama`

**New addon categories:**
- `ai` — adds `internal/ai/client.go`; values: `openai`, `langchaingo`, `gemini`, `ollama`
- `vectorstore` — adds `internal/vectorstore/store.go`; values: `pgvector`, `chromem`, `qdrant`

**`ai-agent` project layout:**
```
<name>/
├── main.go                    # calls agent.Run()
├── agent/
│   └── agent.go               # ReAct-style LLM loop
├── tools/
│   └── tools.go               # sample tool stubs (function calling)
├── llm/
│   └── client.go              # LLM provider initialisation
├── internal/
│   ├── logger/logger.go       # if "other" logging addon selected
│   └── vectorstore/store.go   # if "vectorstore" addon selected
├── go.mod
├── .gitignore
├── Makefile
├── README.md
└── Dockerfile                 # if dockerSupport is true
```

### Phase A — Registry & Metadata

- [x] **T-AI1 — [Backend] Extend `DependencyMap`** — add versioned entries in `registry.go` for: `langchaingo` (`github.com/tmc/langchaingo`), `openai` (`github.com/openai/openai-go`), `gemini` (`github.com/google/generative-ai-go`), `ollama` (`github.com/ollama/ollama`), `pgvector` (`github.com/pgvector/pgvector-go`), `chromem` (`github.com/philippgille/chromem-go`), `qdrant` (`github.com/qdrant/go-client`)

- [x] **T-AI2 — [Backend] Extend Supported\* maps and registries** — in `registry.go`:
  - `SupportedProjectTypesMap` → add `"ai-agent": true`
  - `SupportedProjectTypesLabelsMap` → add `"ai-agent": "AI Agent"`
  - `SupportedFrameworksMap` → add `"ai-agent"` key with `langchaingo | openai | gemini | ollama`
  - `SupportedAddonsMap` → add `"ai"` category (`openai`, `langchaingo`, `gemini`, `ollama`) and `"vectorstore"` category (`pgvector`, `chromem`, `qdrant`)
  - `addonRegistry` → register `&AIAddonGen{}` under `"ai"` and `&VectorStoreAddonGen{}` under `"vectorstore"`
  - `GeneratorRegistry` → register `&AIAgentGenerator{}` under `"ai-agent"`

- [x] **T-AI3 — [Backend] Update `CreateProjectRequest` validate tag** — add `ai-agent` to the `oneof` enum in `types.go`: `oneof=microservice simple-project cli-app api-server ai-agent`

### Phase B — Code Generation Helpers
> Parallel with each other; depends on Phase A.

- [ ] **T-AI4 — [Backend] Add `GenerateLLMClient(framework string) []byte`** — returns `llm/client.go`; provider-specific initialisation:
  - `langchaingo` — `llms.New()` with provider selection via env var
  - `openai` — `openai.NewClient(os.Getenv("OPENAI_API_KEY"))`
  - `gemini` — `genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))`
  - `ollama` — HTTP client pointing at `OLLAMA_HOST` (default `http://localhost:11434`)

- [ ] **T-AI5 — [Backend] Add `GenerateAgentContent(framework, name string) []byte`** — returns `agent/agent.go`; ReAct-style loop per provider: langchaingo agent executor, openai tool-calling loop, gemini function-calling loop, ollama prompt loop with JSON tool parsing

- [ ] **T-AI6 — [Backend] Add `GenerateToolsContent(framework string) []byte`** — returns `tools/tools.go`; sample tool stub using the framework's tool interface (`tools.Tool` for langchaingo, `openai.ChatCompletionToolParam` for openai, `genai.FunctionDeclaration` for gemini, JSON schema struct for ollama)

- [ ] **T-AI7 — [Backend] Add `GenerateAIAddonContent(provider string) []byte`** — returns `internal/ai/client.go`; thin LLM client wrapper for use when the `ai` addon is selected on any non-`ai-agent` project type

- [ ] **T-AI8 — [Backend] Add `GenerateVectorStoreContent(store string) []byte`** — returns `internal/vectorstore/store.go`; store-specific setup:
  - `pgvector` — pgx v5 + pgvector-go connection + sample insert/query
  - `chromem` — `chromem.NewPersistentDB` with sample collection create/query
  - `qdrant` — qdrant gRPC client init + sample upsert/search

All five helpers added to `internal/generator/gen_utils.go`.

### Phase C — Generator & Addon Structs
> Depends on Phase B.

- [ ] **T-AI9 — [Backend] Create `gen_ai_agent.go`** — `AIAgentGenerator.Generate()` wires all Phase B helpers; always emits `main.go`, `agent/agent.go`, `tools/tools.go`, `llm/client.go`, `go.mod` (via `GenerateGoModV2`), `.gitignore`, `Makefile`, `README.md`; handles `"other"` logging addon inline (same pattern as other generators); handles `"vectorstore"` addon via `GenerateVectorStoreContent`; emits `Dockerfile` when `DockerSupport` is true

- [ ] **T-AI10 — [Backend] Add `AIAddonGen` struct** — implements `AddonGenerator`; called when the `ai` addon is selected on any project type; writes `GenerateAIAddonContent(provider)` output to `internal/ai/client.go` in the zip

- [ ] **T-AI11 — [Backend] Add `VectorStoreAddonGen` struct** — implements `AddonGenerator`; called when the `vectorstore` addon is selected on any project type; writes `GenerateVectorStoreContent(store)` output to `internal/vectorstore/store.go` in the zip

`AIAddonGen` and `VectorStoreAddonGen` can be defined inline in `registry.go` or in a dedicated `gen_ai_addon.go`.

---

## Phase 13 — Distribution & Release
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

- [x] **T36 — [Infra] Add `goreleaser.yaml`** — cross-compile `goini` for `darwin/amd64`, `darwin/arm64`, `linux/amd64`, `linux/arm64`; attach binaries to GitHub Releases; include shell completion scripts in archives; configure `brews` block to auto-update the Homebrew formula on each release
- [x] **T37 — [Infra] Update `.github/workflows/release.yml`** — run `goreleaser release` on GitHub Release event; preserve the existing Docker Hub publish job for the server image (`cmd/server`)
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
| 11 — Frontend UI/UX Overhaul | T-UX1–T-UX10 | Medium-High | Phase 1 (stable API) | Open |
| 12 — AI Agent Capabilities | T-AI1–T-AI11 | High | Phase 2 | Open |
| 13 — Distribution & Release | T36–T39 | Low-Medium | Phases 11 + 12 | Open |

---

## Dependencies to Add

### Go (backend)

| Package | Version | Purpose |
|---|---|---|
| `github.com/spf13/cobra` | `v1.10.2` | CLI command framework for `goini` |
| `github.com/charmbracelet/huh` | `latest` | Terminal form / interactive prompts |
| `github.com/tmc/langchaingo` | `latest` | LangChain-style LLM orchestration |
| `github.com/openai/openai-go` | `latest` | Official OpenAI Go SDK |
| `github.com/google/generative-ai-go` | `latest` | Google Gemini SDK |
| `github.com/ollama/ollama` | `latest` | Ollama local LLM client |
| `github.com/pgvector/pgvector-go` | `latest` | pgvector support for Go (pgx v5) |
| `github.com/philippgille/chromem-go` | `latest` | Embedded vector DB (no external service) |
| `github.com/qdrant/go-client` | `latest` | Qdrant gRPC client |

All Go packages have no CGO requirements — the binary remains fully static and cross-compilable.

### Frontend (npm)

| Package | Purpose |
|---|---|
| `vite` + `@vitejs/plugin-react` | Replaces CRA as build tool |
| `@fontsource/geist` | Geist font (self-hosted, no CDN) |
| `react-syntax-highlighter` | Fenced code-block highlighting in Explore docs |

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
| 7. Migrate CRA → Vite (T-UX1) | Phase 11 | All other UX tasks |
| 8. Design token layer + Geist font (T-UX2) | Phase 11 | T-UX3 through T-UX10 |
| 9. App shell + form + docs redesign (T-UX3–T-UX6) | Phase 11 | T-UX7–T-UX10 |
| 10. Responsive, a11y, polish (T-UX7–T-UX10) | Phase 11 | Phase 13 |
| 11. Extend `DependencyMap` + Supported\* maps (T-AI1–T-AI3) | Phase 12 | T-AI4 onward |
| 12. LLM helper generators (T-AI4–T-AI8) | Phase 12 | T-AI9 |
| 13. `AIAgentGenerator` + addon structs (T-AI9–T-AI11) | Phase 12 | Phase 13 |
| 14. goreleaser + release workflow + Homebrew + README (T36–T39) | Phase 13 | — |
