# go-initializer

go-initializer scaffolds production-ready Go projects in seconds — from a terminal or a browser. Pick a project type, framework, and addons, and get a fully wired, immediately runnable codebase dropped into your working directory.

---

## Installation

### CLI (recommended)

```sh
go install github.com/neo7337/go-initializer/cmd/goini@latest
```

### Homebrew

```sh
brew tap neo7337/goini
brew install goini
```

### Web UI

The hosted web UI is available at **[https://go-initializer.dev](https://go-initializer.dev)** — no installation required.

---

## Quick Start

### Interactive mode

```sh
goini new
```

The CLI walks you through every option step-by-step. Any flag you provide skips its corresponding prompt, making the command fully scriptable.

### Scripted / CI mode

```sh
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

The project is extracted directly into the output directory — no zip file left on disk.

---

## Commands

```
goini new                       scaffold a new project (interactive or via flags)
goini list types                print all supported project types
goini list frameworks           print frameworks for a type  (--type <t>)
goini list addons               print all supported addon categories and values
goini version                   print version and build info
goini completion bash|zsh|fish  print shell completion script
```

---

## Supported Options

### Go versions

`1.25.0` · `1.24.6` · `1.23.12`

### Project types & frameworks

| Project Type   | Frameworks                              |
|----------------|-----------------------------------------|
| Microservice   | Gin · Echo · Fiber · Golly · Go kit     |
| API Server     | Gin · Echo · Fiber · Chi · Golly        |
| CLI App        | Cobra · urfave/cli · Kingpin · Golly    |
| Simple Project | Golly                                   |
| AI Agent       | LangChainGo · OpenAI · Gemini · Ollama  |

### Addons

| Category     | Values                                      |
|--------------|---------------------------------------------|
| `cache`      | `redis` · `memcached`                       |
| `database`   | `gorm` · `ent`                              |
| `other`      | `zap` · `logrus` · `cobra`                  |
| `ai`         | `openai` · `langchaingo` · `gemini` · `ollama` |
| `vectorstore`| `pgvector` · `chromem` · `qdrant`           |

---

## What gets generated

Every project includes:

- `go.mod` wired with all selected dependencies
- Framework-specific `main.go`, handlers, and router
- `.gitignore` (binaries, `vendor/`, `.env`, IDE dirs)
- `Makefile` with `build`, `run`, `tidy`, and `test` targets
- `README.md` with next-step instructions
- Optional `Dockerfile` (multi-stage) when `--docker` is set
- Optional addon files: `internal/cache/`, `internal/database/`, `internal/logger/`, `internal/ai/`, `internal/vectorstore/`

### Next steps after generation

```sh
cd <name>
go mod tidy
make run      # server projects
make build    # CLI projects
```

---

## Project Structure

```
go-initializer/
├── cmd/
│   ├── goini/          # CLI binary
│   └── server/         # HTTP server binary
├── internal/
│   ├── generator/      # shared generation engine
│   └── server/         # HTTP layer
└── frontend/           # React web UI
```

---

## Docker (self-hosted web UI)

```sh
docker-compose up --build
```

The backend listens on `:8182` and the frontend on `:8001`.

---

## Contributing

Contributions are welcome! Please open issues or pull requests for suggestions and improvements. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT
