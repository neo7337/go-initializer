# AI Agent Guide

Go Initializer's **AI Agent** project type scaffolds a production-ready Go agent wired to your chosen LLM provider. This guide walks through the generated structure, configuration, workflow integration patterns, and production deployment.

## Overview

When you select **AI Agent** in the generator:

- The **Framework** picker becomes an **LLM Provider** selector (`langchaingo`, `openai`, `gemini`, `ollama`)
- The generated project ships a complete agent loop with tool-calling out of the box
- Optionally add a **vector store** add-on (`pgvector`, `chromem`, `qdrant`) for Retrieval-Augmented Generation

---

## Generated project layout

```
my-agent/
├── go.mod
├── main.go            # entry point — initialises the agent and starts execution
├── agent/
│   └── agent.go       # core agent loop with tool dispatch
├── tools/
│   └── tools.go       # tool definitions (extend here)
├── llm/
│   └── client.go      # LLM provider client (provider-specific)
├── internal/
│   └── vectorstore/   # present when a vector store add-on is selected
│       └── store.go
├── Makefile
├── README.md
└── Dockerfile         # multi-stage build (if Docker support enabled)
```

---

## Configuring the LLM provider

### LangChainGo

Set your provider key in the environment. The generated `llm/client.go` reads it automatically:

```bash
export OPENAI_API_KEY=sk-...
go run .
```

To switch provider backend, update the import in `llm/client.go`:

```go
// OpenAI (default)
import "github.com/tmc/langchaingo/llms/openai"

// Google Gemini
import "github.com/tmc/langchaingo/llms/googleai"

// Ollama (local)
import "github.com/tmc/langchaingo/llms/ollama"
```

No other code needs to change — all providers satisfy the same `llms.Model` interface.

### OpenAI

```bash
export OPENAI_API_KEY=sk-...
go run .
```

The client uses `github.com/openai/openai-go` and defaults to `gpt-4o`. Change the model in `llm/client.go`:

```go
resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
    Model: openai.ChatModelGPT4o,
    ...
})
```

### Gemini

```bash
export GEMINI_API_KEY=AIza...
go run .
```

### Ollama (local, no API key)

```bash
# Start Ollama with a model
ollama pull llama3
ollama serve

go run .   # connects to localhost:11434 by default
```

Override the host:

```bash
export OLLAMA_HOST=http://my-ollama-server:11434
```

---

## Adding tools

Tools are the heart of an agent — they let the LLM take actions in the world. Open `tools/tools.go` and add your own:

```go
// tools/tools.go
package tools

import "github.com/tmc/langchaingo/tools"

// SearchWeb is an example custom tool.
type SearchWeb struct{}

func (s SearchWeb) Name() string        { return "search_web" }
func (s SearchWeb) Description() string { return "Search the internet for up-to-date information." }

func (s SearchWeb) Call(ctx context.Context, input string) (string, error) {
    // implement your search logic here
    return "result from web search", nil
}
```

Register it in `agent/agent.go`:

```go
agentTools := []tools.Tool{
    tools.Calculator{},
    tools.SearchWeb{},   // ← add here
}
```

---

## Vector store / RAG

When a vector store add-on is selected, `internal/vectorstore/store.go` is generated with a ready-to-use client. Here's the typical RAG pattern:

```go
// 1. Embed and index documents
docs := []schema.Document{
    {PageContent: "Go 1.22 adds range-over-integer loops.", Metadata: map[string]any{"src": "release-notes"}},
}
if err := store.AddDocuments(ctx, docs); err != nil {
    log.Fatal(err)
}

// 2. Retrieve relevant context at query time
results, err := store.SimilaritySearch(ctx, query, 3)
if err != nil {
    log.Fatal(err)
}

// 3. Inject context into the LLM prompt
context := buildContext(results)
```

See [Working with Add-ons](how-to-addons) for vector store-specific setup (Postgres extension, Qdrant Docker, etc.).

---

## Workflow integration

The generated agent is a self-contained Go binary — it can be embedded into any orchestration system.

### As a standalone service

Wrap the agent in an HTTP handler to expose it as a microservice:

```go
http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
    var req struct{ Input string }
    json.NewDecoder(r.Body).Decode(&req)

    result, err := myAgent.Run(r.Context(), req.Input)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"output": result})
})
http.ListenAndServe(":8080", nil)
```

### With a task queue (e.g. Temporal, Asynq)

The agent `Run()` method is just a function call — wrap it in a worker:

```go
// Temporal activity
func RunAgentActivity(ctx workflow.Context, input string) (string, error) {
    return myAgent.Run(ctx, input)
}

// Asynq task handler
func HandleAgentTask(ctx context.Context, t *asynq.Task) error {
    result, err := myAgent.Run(ctx, string(t.Payload()))
    // …persist result
    return err
}
```

### With n8n / Zapier / Make via webhook

Deploy the agent as an HTTP service (see above), then call it from any no-code automation platform using a **Webhook** / **HTTP Request** node pointing at your deployed URL.

### In a LangChain pipeline

Because the LangChainGo provider implements `llms.Model`, the agent can be composed inside larger chains:

```go
chain := chains.NewConversation(llm, memory.NewConversationBuffer())
// or used as a sub-agent in a router chain
```

---

## Docker deployment

If **Docker support** was enabled, the project ships a multi-stage `Dockerfile`:

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o agent .

# Runtime stage
FROM alpine:latest
COPY --from=builder /app/agent /agent
ENV OPENAI_API_KEY=""
ENTRYPOINT ["/agent"]
```

Build and run:

```bash
docker build -t my-agent .
docker run -e OPENAI_API_KEY=sk-... my-agent
```

---

## Environment variable reference

| Variable | Provider | Description |
|----------|----------|-------------|
| `OPENAI_API_KEY` | OpenAI, LangChainGo (OpenAI backend) | OpenAI API key |
| `GEMINI_API_KEY` | Gemini, LangChainGo (Google backend) | Google AI API key |
| `OLLAMA_HOST` | Ollama | Ollama server URL (default: `http://localhost:11434`) |
| `DATABASE_URL` | pgvector add-on | PostgreSQL DSN (`postgres://user:pass@host/db`) |
| `QDRANT_URL` | Qdrant add-on | Qdrant gRPC URL (default: `http://localhost:6334`) |

---

## Next steps

- [Working with Add-ons](how-to-addons) — set up pgvector, Qdrant, or chromem
- [Choosing a Framework](frameworks) — compare LLM providers in depth
- [REST API Reference](api-reference) — generate an AI Agent project via cURL
