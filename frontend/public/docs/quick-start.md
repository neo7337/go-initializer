# Quick Start

Get a Go project generated and running in under two minutes.

## 1. Choose a project type

Open the generator and pick one of the five project types:

| Type | Best for |
|------|----------|
| **Simple Project** | Minimal single-binary HTTP service |
| **Microservice** | Service in a broader microservices system |
| **API Server** | Public REST API with middleware |
| **CLI Application** | Terminal tools and dev utilities |
| **AI Agent** | LLM-backed agent, RAG pipeline, or workflow integration |

## 2. Fill in the details

| Field | Description |
|-------|-------------|
| Module Name | Your Go module path, e.g. `github.com/you/myapp` |
| Project Name | The name of the root directory |
| Go Version | Select the Go runtime version |
| Framework | HTTP framework (Gin, Echo, Fiber, …) — or LLM provider for AI Agent |

> **AI Agent note:** When **AI Agent** is selected the Framework picker shows LLM providers instead of HTTP frameworks: `langchaingo`, `openai`, `gemini`, `ollama`.

## 3. Pick add-ons (optional)

Enable any combination of:

- **Databases** — GORM or Ent ORM wiring  
- **Cache** — Redis or Memcached client setup  
- **Logging** — Zap or Logrus structured logger  
- **AI** — thin LLM client wrapper (available on all types)  
- **Vector Store** — pgvector, Qdrant, or chromem for RAG workflows  
- **Docker** — `Dockerfile` with multi-stage build  

## 4. Generate & download

Click **Generate Project** — a `.zip` archive downloads immediately. Unzip and run:

```bash
cd your-project
go mod tidy
go run .
```

Your server is live on `localhost:8080`.

---

## AI Agent quick start

```bash
# Generate
# Type: AI Agent · Provider: langchaingo · Add-on: vectorstore=pgvector

cd my-agent
go mod tidy

# Set your provider key
export OPENAI_API_KEY=sk-...   # or GEMINI_API_KEY / OLLAMA_HOST

go run .
```

The generated `agent/agent.go` exposes a ready-to-run agent loop. Drop in your own tools in `tools/tools.go` and replace the placeholder prompt to fit your use-case.

See the [AI Agent Guide](ai-agent) for production patterns and workflow integration.

