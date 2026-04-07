<div align="center">

# go-initializer

**Scaffold any Go project — microservice, API, CLI, or AI agent — in seconds.**

Pick a project type, framework, and addons. Get a fully wired, immediately runnable codebase dropped into your working directory.

[![Backend CI](https://github.com/neo7337/go-initializer/actions/workflows/backend.yml/badge.svg)](https://github.com/neo7337/go-initializer/actions/workflows/backend.yml)
[![Frontend CI](https://github.com/neo7337/go-initializer/actions/workflows/frontend.yml/badge.svg)](https://github.com/neo7337/go-initializer/actions/workflows/frontend.yml)
[![Security Scan](https://github.com/neo7337/go-initializer/actions/workflows/security.yml/badge.svg)](https://github.com/neo7337/go-initializer/actions/workflows/security.yml)
[![Latest Release](https://img.shields.io/github/v/release/neo7337/go-initializer?color=blue)](https://github.com/neo7337/go-initializer/releases/latest)
[![Go Version](https://img.shields.io/badge/go-1.25-00ADD8?logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/neo7337/goinitializer?label=docker%20pulls)](https://hub.docker.com/r/neo7337/goinitializer)

[Web UI](https://go-initializer.dev) · [Docs](https://go-initializer.dev/docs) · [Contributing](CONTRIBUTING.md) · [Roadmap](https://github.com/neo7337/go-initializer/projects)

</div>

---

## What is go-initializer?

go-initializer is a Go project scaffolding tool — the Go ecosystem's answer to Spring Initializr or `create-react-app`. In one command (or a few browser clicks), you get:

- A `go.mod` wired with every dependency you selected
- Framework-specific `main.go`, handlers, and router
- A `Makefile` with `build`, `run`, `tidy`, and `test` targets
- A multi-stage `Dockerfile` (optional)
- A `.gitignore`, a starter `README.md`, and addon packages under `internal/`

**No manual boilerplate. No forgotten `go mod tidy`. Just code.**

---

## Install

### Homebrew (macOS / Linux)

```sh
brew tap neo7337/goini
brew install goini
```

### go install

```sh
go install github.com/neo7337/go-initializer/cmd/goini@latest
```

### Web UI — no install required

Open **[https://go-initializer.dev](https://go-initializer.dev)** in your browser and download a zip.

---

## Quick start

### Interactive mode (recommended for first-timers)

```sh
goini new
```

The CLI walks you through every option step-by-step via an interactive TUI.

### Scripted / CI mode

Any flag you provide skips its corresponding prompt — making the command fully scriptable:

```sh
goini new \
  --name        myapp \
  --module      github.com/acme/myapp \
  --description "A production microservice" \
  --go-version  1.25.0 \
  --type        microservice \
  --framework   gin \
  --addon       cache=redis \
  --addon       database=gorm \
  --addon       other=zap \
  --docker \
  --output      ./projects/
```

The project is extracted directly into the output directory — no zip file left on disk.

### REST API

```sh
curl -s -X POST https://go-initializer.dev/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "name": "myapp",
    "module": "github.com/acme/myapp",
    "goVersion": "1.25.0",
    "projectType": "microservice",
    "framework": "gin",
    "addons": [{"category":"cache","value":"redis"}],
    "docker": true
  }' --output myapp.zip
```

---

## Supported options

### Go versions

`1.25.0` · `1.24.6` · `1.23.12`

### Project types & frameworks

| Project type   | Frameworks                                    |
|----------------|-----------------------------------------------|
| Microservice   | Gin · Echo · Fiber · Golly · Go kit           |
| API Server     | Gin · Echo · Fiber · Chi · Golly              |
| CLI App        | Cobra · urfave/cli · Kingpin · Golly          |
| Simple Project | Golly                                         |
| AI Agent       | LangChainGo · OpenAI · Gemini · Ollama        |

### Addons

| Category      | Values                                              |
|---------------|-----------------------------------------------------|
| `cache`       | `redis` · `memcached`                               |
| `database`    | `gorm` · `ent`                                      |
| `other`       | `zap` · `logrus` · `cobra`                          |
| `ai`          | `openai` · `langchaingo` · `gemini` · `ollama`      |
| `vectorstore` | `pgvector` · `chromem` · `qdrant`                   |

---

## Why go-initializer?

| | go-initializer | Manual setup | Other generators |
|---|---|---|---|
| Interactive TUI | ✅ | ❌ | Varies |
| Web UI (zero install) | ✅ | ❌ | Rare |
| REST API | ✅ | ❌ | Rare |
| AI Agent scaffolding | ✅ | ❌ | ❌ |
| Vector store wiring | ✅ | ❌ | ❌ |
| Multi-stage Dockerfile | ✅ | Manual | Rare |
| Makefile included | ✅ | Manual | Rare |
| In-memory generation (no temp files) | ✅ | N/A | Varies |
| CI/scripted flag mode | ✅ | N/A | Rare |
| Self-hostable | ✅ | N/A | Rare |

---

## What gets generated

Every project ships with:

- `go.mod` — module declaration wired with all selected dependencies
- `main.go` — framework-specific entry point, handlers, and router
- `.gitignore` — binaries, `vendor/`, `.env`, IDE directories
- `Makefile` — `build`, `run`, `tidy`, `test` targets ready to use
- `README.md` — next-step instructions tailored to the generated project
- `Dockerfile` *(optional)* — multi-stage build with minimal final image
- `internal/cache/` *(optional)* — Redis or Memcached client wiring
- `internal/database/` *(optional)* — GORM or Ent client wiring
- `internal/logger/` *(optional)* — Zap or Logrus setup
- `internal/ai/` *(optional)* — LLM provider client + tool-calling scaffold
- `internal/vectorstore/` *(optional)* — pgvector, Qdrant, or chromem wiring

### Next steps after generation

```sh
cd <name>
go mod tidy
make run      # server / AI agent projects
make build    # CLI projects
```

---

## CLI reference

```
goini new                       scaffold a new project (interactive or via flags)
goini list types                print all supported project types
goini list frameworks           print frameworks for a type  (--type <t>)
goini list addons               print all supported addon categories and values
goini version                   print version and build info
goini completion bash|zsh|fish  print shell completion script
```

---

## Self-hosting

Run the full stack locally with Docker Compose:

```sh
docker-compose up --build
```

| Service  | Default port |
|----------|-------------|
| Backend  | `:8182`     |
| Frontend | `:8001`     |

See [docs/how-to-self-host](https://go-initializer.dev/docs/how-to-self-host) for environment variable reference and production deployment notes.

---

## Project structure

```
go-initializer/
├── cmd/
│   ├── goini/          # CLI binary (cobra + charmbracelet huh TUI)
│   └── server/         # HTTP server binary (Gin)
├── internal/
│   ├── generator/      # shared generation engine (used by both CLI and server)
│   └── server/         # HTTP layer, middleware, rate limiting
└── frontend/           # React + TypeScript web UI (Vite + Radix UI + Tailwind)
```

The CLI and web UI share the exact same `internal/generator` package — output is guaranteed identical regardless of which interface you use.

---

## Contributing

Contributions are welcome — bug reports, feature requests, documentation improvements, and code PRs alike.

- Read [CONTRIBUTING.md](CONTRIBUTING.md) for the dev setup and workflow
- Browse [`good first issue`](https://github.com/neo7337/go-initializer/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) for beginner-friendly tasks
- Open a [Discussion](https://github.com/neo7337/go-initializer/discussions) for questions or ideas

### Contributors

Thanks to everyone who has contributed to go-initializer:

[![Contributors](https://contrib.rocks/image?repo=neo7337/go-initializer)](https://github.com/neo7337/go-initializer/graphs/contributors)

---

## License

[MIT](LICENSE)
