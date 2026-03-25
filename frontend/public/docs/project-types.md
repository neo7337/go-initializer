# Project Types

Go Initializer supports five project types, each producing a different scaffold structure.

## Simple Project

A minimal single-binary Go service backed by **Golly** (the only supported framework for this type). Use this as a clean starting point when you want full control and no heavy third-party routing layer.

**Generated layout:**

```
my-project/
├── go.mod
├── main.go
├── handler.go
└── Dockerfile          # if Docker support is enabled
```

**When to choose:** Internal tooling, one-off microservices, or learning projects where you want the least amount of scaffolding.

---

## Microservice

A production-oriented layout for services that will run behind a gateway or alongside other services. Supports **Gin**, **Echo**, **Fiber**, **Golly**, and **GoKit**.

**Generated layout:**

```
my-service/
├── go.mod
├── main.go
├── internal/
│   ├── handler/
│   │   └── handler.go
│   ├── router/
│   │   └── router.go
│   └── server/
│       └── server.go
├── pkg/
└── Dockerfile
```

Health endpoints (`/healthz`, `/readyz`) are wired automatically.

**When to choose:** Services that are part of a larger system and need structured internals, health probes, and a clean handler/router split.

---

## API Server

Identical intent to Microservice but targets teams that prefer the `api-server` naming convention and want Chi or a more stdlib-compatible routing style. Supports **Gin**, **Echo**, **Fiber**, **Chi**, and **Golly**.

**Generated layout:**

```
my-api/
├── go.mod
├── main.go
├── internal/
│   ├── handler/
│   ├── middleware/
│   └── server/
└── Dockerfile
```

**When to choose:** REST APIs, BFF (backend-for-frontend) services, or any HTTP layer that benefits from middleware composability.

---

## CLI Application

Scaffolds a command-line binary with a proper sub-command structure. Supports **Cobra**, **Urfave/cli**, **Kingpin**, and **Golly**.

**Generated layout:**

```
my-cli/
├── go.mod
├── main.go
└── cmd/
    ├── root.go
    └── version.go
```

**When to choose:** DevOps tooling, database migration runners, developer utilities, or any program that is invoked from a terminal rather than over HTTP.

---

## AI Agent

Scaffolds a production-ready AI agent backed by your chosen LLM provider. The `framework` field selects the provider: **LangChainGo**, **OpenAI**, **Gemini**, or **Ollama**. Optionally wire in a vector store (pgvector, Qdrant, or chromem) via the add-ons panel.

**Generated layout:**

```
my-agent/
├── go.mod
├── main.go
├── agent/
│   └── agent.go        # core agent loop with tool-calling
├── tools/
│   └── tools.go        # tool definitions
├── llm/
│   └── client.go       # LLM provider client
├── Makefile
├── README.md
└── Dockerfile          # if Docker support is enabled
```

> **Note:** For AI Agent projects, the "Framework" selector in the generator maps to your LLM provider (LangChainGo, OpenAI, Gemini, Ollama) rather than an HTTP framework.

**When to choose:** Autonomous agents, Retrieval-Augmented Generation (RAG) pipelines, LLM-backed tooling, or any project that needs to integrate with an AI/ML workflow or orchestration system.

See the [AI Agent Guide](ai-agent) for a full walkthrough including vector store integration and workflow patterns.

---

## Choosing the right type

| Situation | Recommended type |
| --- | --- |
| Simple HTTP service, no structure needed | Simple Project |
| Service in a microservices mesh | Microservice |
| Public-facing REST API | API Server |
| Command-line tool | CLI Application |
| LLM-backed agent or RAG pipeline | AI Agent |

## Simple Project

A minimal single-binary Go service backed by **Golly** (the only supported framework for this type). Use this as a clean starting point when you want full control and no heavy third-party routing layer.

**Generated layout:**

```
my-project/
├── go.mod
├── main.go
├── handler.go
└── Dockerfile          # if Docker support is enabled
```

**When to choose:** Internal tooling, one-off microservices, or learning projects where you want the least amount of scaffolding.

---

## Microservice

A production-oriented layout for services that will run behind a gateway or alongside other services. Supports **Gin**, **Echo**, **Fiber**, **Golly**, and **GoKit**.

**Generated layout:**

```
my-service/
├── go.mod
├── main.go
├── internal/
│   ├── handler/
│   │   └── handler.go
│   ├── router/
│   │   └── router.go
│   └── server/
│       └── server.go
├── pkg/
└── Dockerfile
```

Health endpoints (`/healthz`, `/readyz`) are wired automatically.

**When to choose:** Services that are part of a larger system and need structured internals, health probes, and a clean handler/router split.

---

## API Server

Identical intent to Microservice but targets teams that prefer the `api-server` naming convention and want Chi or a more stdlib-compatible routing style. Supports **Gin**, **Echo**, **Fiber**, **Chi**, and **Golly**.

**Generated layout:**

```
my-api/
├── go.mod
├── main.go
├── internal/
│   ├── handler/
│   ├── middleware/
│   └── server/
└── Dockerfile
```

**When to choose:** REST APIs, BFF (backend-for-frontend) services, or any HTTP layer that benefits from middleware composability.

---

## CLI Application

Scaffolds a command-line binary with a proper sub-command structure. Supports **Cobra**, **Urfave/cli**, **Kingpin**, and **Golly**.

**Generated layout:**

```
my-cli/
├── go.mod
├── main.go
└── cmd/
    ├── root.go
    └── version.go
```

**When to choose:** DevOps tooling, database migration runners, developer utilities, or any program that is invoked from a terminal rather than over HTTP.

---

## Choosing the right type

| Situation | Recommended type |
| --- | --- |
| Simple HTTP service, no structure needed | Simple Project |
| Service in a microservices mesh | Microservice |
| Public-facing REST API | API Server |
| Command-line tool | CLI Application |
