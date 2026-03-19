# Application Architecture

Go Initializer is composed of three layers that work together to scaffold Go projects on demand.

## Overview

```
┌─────────────────────────────────────────────────────────┐
│                      Browser / CLI                       │
│              React UI  ·  goini binary                   │
└──────────────────────┬──────────────────────────────────┘
                       │  HTTP (REST)
┌──────────────────────▼──────────────────────────────────┐
│                    Go HTTP Server                        │
│         Gin router  ·  request validation               │
│         GET /api/meta   ·   POST /api/generate          │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                 Generator Engine                         │
│   GeneratorRegistry  →  per-type Generator              │
│   AddonRegistry      →  per-category AddonGenerator     │
│   In-memory zip assembly  →  []byte response            │
└─────────────────────────────────────────────────────────┘
```

## Frontend

A React + TypeScript single-page application built with Vite. It fetches the available options from `/api/meta` on load, presents the form, and POSTs to `/api/generate` to receive the zip archive.

The UI is entirely stateless — it holds no session and makes no assumptions about the user; every interaction produces a deterministic HTTP request.

## Backend

A Go HTTP server powered by the **Gin** framework. Two endpoints handle all traffic:

| Endpoint | Method | Purpose |
| --- | --- | --- |
| `/api/meta` | GET | Returns supported Go versions, project types, frameworks, and addons |
| `/api/generate` | POST | Validates the request and returns a `.zip` archive |
| `/healthz` | GET | Liveness probe — returns 200 when the process is alive |

The server performs structural validation with `go-playground/validator` before the request ever reaches the generator.

## Generator Engine

The engine lives in `internal/generator/`. Each project type registers itself in `GeneratorRegistry`:

```go
var GeneratorRegistry = map[string]Generator{
    "simple-project": &SimpleProjectGenerator{},
    "microservice":   &MicroserviceGenerator{},
    "cli-app":        &CLIAppGenerator{},
}
```

When a request arrives, the matching `Generator.Generate()` is called. It assembles files directly into a `zip.Writer` backed by an in-memory `bytes.Buffer` — no temp files, no disk I/O.

Addons are handled by a second registry (`addonRegistry`) that writes additional files — database models, cache clients, logger setup — alongside the core scaffold.

## goini CLI

The `goini` binary at `cmd/goini/` is an alternative entry point. It either runs an interactive TUI powered by **Charmbracelet Huh** or accepts flags for fully non-interactive (CI-friendly) use. Internally it calls the same generator engine, so the output is identical to the web UI.

## Data Flow (generate request)

1. User fills in the form (or runs `goini new`)
2. A `CreateProjectRequest` JSON is sent to `POST /api/generate`
3. The server validates required fields and allowed values
4. The matched generator builds a zip in memory
5. Addon generators append extra files to the zip
6. The zip bytes are streamed back as `application/zip`
7. The browser triggers a `project.zip` download
