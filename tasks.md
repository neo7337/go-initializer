# go-initializer — Implementation Task List

## Goal
The backend receives form input (project type, Go version, framework, addons, docker support, module name, project name, description) and generates a downloadable `.zip` of a complete, immediately runnable Go project. Tasks are ordered **low → high complexity** within each phase so implementation can proceed step by step without blockers.

---

## Current State Overview

### Backend
| File | Status |
|---|---|
| Server, router, config | Complete |
| `GET /api/meta` | Complete |
| `POST /api/generate` (routing) | Complete — registry-based, no switch statement |
| `registry.go` — `Generator` / `AddonGenerator` interfaces + registries | **New** — Complete |
| `gen_utils.go` — `addToZip` helper, `CacheAddonGen`, `DatabaseAddonGen` | **New** — Complete |
| `gen_utils.go` — `GenerateGoModV2`, `GenerateDockerfile` | Working |
| `gen_utils.go` — `GenerateCacheAddon` | Working |
| `gen_utils.go` — `GenerateDatabaseAddon` | Fixed — uses `os.Getenv("DATABASE_URL")`, no crashes |
| `gen_utils.go` — `GenerateMainContent` | **Hardcoded to golly/l3** — not framework-aware |
| `gen_simple_project.go` — `SimpleProjectGenerator` | Refactored — uses `addToZip` + addon registry; still calls framework-unaware `GenerateMainContent` |
| `gen_microservice.go` — `MicroserviceGenerator` | Substantially implemented — generates README, go.mod, cmd/main.go, internal/handler + service, addons, Dockerfile. **main.go is not framework-aware** |
| `router.go` — CORS | Fixed — removed invalid `AllowCredentials: true` with wildcard origin |
| CLI App generator | **Missing entirely** |
| API Server generator | **Missing entirely** |
| `handler.go` | Uses `generatorRegistry` — adding new types requires only a registry entry, no handler edits |
| Input validation | `validator` instance created in `server.go` but never used; no struct tags on `CreateProjectRequest` |

### Frontend
| File | Status |
|---|---|
| UI layout, form, theming | Complete |
| `useGeneratorForm.ts` hook | Complete |
| `service.ts` — `getMetaData` | Complete |
| `service.ts` — `generateProject` | Fixed — correctly calls `/api/generate` |
| `GeneratorForm.tsx` — addon options | **Hardcoded** — not driven by API metadata |
| `GeneratorForm.tsx` — `logrus` addon | Missing (backend supports it, UI doesn't show it) |
| Error handling | Uses `alert()` — no inline error UI |
| Post-generate feedback | Silent — no success indicator |

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
> Low complexity. Housekeeping tasks, no logic change.

- [ ] **T24 — [Backend] Resolve dead config code** — `config.go` and `loadConfig` are never called; either wire `config.yaml` + call `loadConfig` in `Start()` to drive host/port/timeouts, or delete the dead code
- [ ] **T25 — [Frontend/Infra] Validate `docker-compose.yml`** — verify backend port `8182` and the nginx upstream proxy in `frontend/nginx.conf` are consistent end-to-end

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
| 8 — Infrastructure | T24–T25 | Low | Any | Open |
