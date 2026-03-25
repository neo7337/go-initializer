# Choosing a Go Web Framework

Go's standard library `net/http` is capable on its own, but most teams reach for a framework for routing, middleware, and ergonomics. Here's a comparison of the frameworks supported by Go Initializer.

> **AI Agent projects** use a different selector — see [LLM Providers](#llm-providers-ai-agent) below.

## Gin

**GitHub Stars:** ~80k · **License:** MIT

The most widely used Go web framework. Fast router, minimal allocations, a large middleware ecosystem.

```go
r := gin.Default()
r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
})
r.Run()
```

**Best for:** Teams that want battle-tested stability and a huge community to draw on.

---

## Echo

**GitHub Stars:** ~30k · **License:** MIT

Minimalist and highly performant. Clean API, first-class middleware support, built-in data binding and validation.

```go
e := echo.New()
e.GET("/", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
})
e.Logger.Fatal(e.Start(":1323"))
```

**Best for:** Teams that prefer a clean, opinionated API with strong documentation.

---

## Fiber

**GitHub Stars:** ~35k · **License:** MIT

Inspired by Express.js. Built on top of `fasthttp` for extreme throughput. Near-zero memory allocations.

```go
app := fiber.New()
app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
})
app.Listen(":3000")
```

**Best for:** High-throughput services where raw performance is critical.

---

## Chi

**GitHub Stars:** ~18k · **License:** MIT

`net/http`-compatible router with composable middleware. No external dependencies.

```go
r := chi.NewRouter()
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("welcome"))
})
http.ListenAndServe(":3000", r)
```

**Best for:** Teams that want to stay close to the stdlib and avoid framework lock-in.

---

## Summary

| Framework | Speed | Ecosystem | Stdlib-compatible |
|-----------|-------|-----------|-------------------|
| Gin       | ★★★★  | ★★★★★     | No                |
| Echo      | ★★★★  | ★★★★      | No                |
| Fiber     | ★★★★★ | ★★★       | No (fasthttp)     |
| Chi       | ★★★★  | ★★★       | Yes               |

---

## LLM Providers (AI Agent)

When **AI Agent** is selected as the project type, the Framework picker becomes an LLM provider selector. The generated project's `llm/client.go` is wired to that provider.

### LangChainGo

The Go port of LangChain. Supports multiple providers under a single abstraction, chains, memory, and tool calling.

```go
import "github.com/tmc/langchaingo/llms/openai"

llm, err := openai.New()
```

**Best for:** Teams coming from Python LangChain, or that want a provider-agnostic abstraction with built-in chains and tool support.

---

### OpenAI

Direct integration with the OpenAI API via the official Go SDK. Gives fine-grained access to Chat Completions, Assistants, function calling, and embeddings.

```go
import "github.com/openai/openai-go"

client := openai.NewClient() // reads OPENAI_API_KEY from env
```

**Best for:** Projects already standardised on OpenAI models (GPT-4o, o1, etc.) or that need the latest API features immediately.

---

### Gemini

Integration with Google's Gemini models via the `google/generative-ai-go` SDK.

```go
import "github.com/google/generative-ai-go/genai"

client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
```

**Best for:** Teams on Google Cloud, or that need Gemini's multimodal and long-context capabilities.

---

### Ollama

Runs models locally via [Ollama](https://ollama.com). No API key required — models run on your own hardware.

```go
import "github.com/ollama/ollama/api"

client, err := api.ClientFromEnvironment() // OLLAMA_HOST defaults to localhost:11434
```

**Best for:** Air-gapped deployments, cost-sensitive projects, privacy-first use-cases, or rapid local experimentation.

---

## Provider comparison

| Provider | Hosted | Local | Tool Calling | Embeddings | Free tier |
|----------|--------|-------|-------------|------------|-----------|
| LangChainGo | ✓ (multi) | ✓ | ✓ | ✓ | Depends on backend |
| OpenAI | ✓ | — | ✓ | ✓ | — |
| Gemini | ✓ | — | ✓ | ✓ | ✓ |
| Ollama | — | ✓ | Model-dependent | ✓ | ✓ (self-hosted) |


## Gin

**GitHub Stars:** ~80k · **License:** MIT

The most widely used Go web framework. Fast router, minimal allocations, a large middleware ecosystem.

```go
r := gin.Default()
r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
})
r.Run()
```

**Best for:** Teams that want battle-tested stability and a huge community to draw on.

---

## Echo

**GitHub Stars:** ~30k · **License:** MIT

Minimalist and highly performant. Clean API, first-class middleware support, built-in data binding and validation.

```go
e := echo.New()
e.GET("/", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
})
e.Logger.Fatal(e.Start(":1323"))
```

**Best for:** Teams that prefer a clean, opinionated API with strong documentation.

---

## Fiber

**GitHub Stars:** ~35k · **License:** MIT

Inspired by Express.js. Built on top of `fasthttp` for extreme throughput. Near-zero memory allocations.

```go
app := fiber.New()
app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
})
app.Listen(":3000")
```

**Best for:** High-throughput services where raw performance is critical.

---

## Chi

**GitHub Stars:** ~18k · **License:** MIT

`net/http`-compatible router with composable middleware. No external dependencies.

```go
r := chi.NewRouter()
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("welcome"))
})
http.ListenAndServe(":3000", r)
```

**Best for:** Teams that want to stay close to the stdlib and avoid framework lock-in.

---

## Summary

| Framework | Speed | Ecosystem | Stdlib-compatible |
|-----------|-------|-----------|-------------------|
| Gin       | ★★★★  | ★★★★★     | No                |
| Echo      | ★★★★  | ★★★★      | No                |
| Fiber     | ★★★★★ | ★★★       | No (fasthttp)     |
| Chi       | ★★★★  | ★★★       | Yes               |
