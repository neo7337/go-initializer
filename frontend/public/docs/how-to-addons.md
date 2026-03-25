# How to Work with Add-ons

Add-ons inject pre-wired integrations into your generated project. This guide explains what each add-on generates and how to configure it after download.

## Selecting add-ons

In the web UI, add-ons appear in Section 04. You can combine add-ons freely across categories.

From the CLI:

```bash
goini new --addon cache=redis --addon database=gorm --addon other=zap
```

---

## Cache

### redis

Adds a `pkg/cache/redis.go` file with a configured `go-redis` client:

```go
import "github.com/redis/go-redis/v9"

func NewRedisClient(addr string) *redis.Client {
    return redis.NewClient(&redis.Options{Addr: addr})
}
```

**After generation:**

```bash
go mod tidy   # pulls in github.com/redis/go-redis/v9
```

Point the client at your Redis instance by passing the address (e.g. `localhost:6379`).

### memcached

Adds a `pkg/cache/memcached.go` file using `bradfitz/gomemcache`:

```go
import "github.com/bradfitz/gomemcache/memcache"

func NewMemcachedClient(addr ...string) *memcache.Client {
    return memcache.New(addr...)
}
```

---

## Database

### gorm

Adds `pkg/database/gorm.go` with a GORM connection helper and a sample model:

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

func NewDB(dsn string) (*gorm.DB, error) {
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

Switch the driver import (`sqlite`, `mysql`, `sqlserver`) to match your database.

### ent

Adds the Ent schema directory (`ent/schema/`) and a client initialiser. Run the Ent code-generator after customising your schema:

```bash
go run entgo.io/ent/cmd/ent generate ./ent/schema
```

---

## Other

### zap

Adds `pkg/logger/zap.go` with a production-ready Zap logger:

```go
import "go.uber.org/zap"

func NewLogger() (*zap.Logger, error) {
    return zap.NewProduction()
}
```

Use `zap.NewDevelopment()` locally for human-readable output.

### logrus

Adds `pkg/logger/logrus.go` with a Logrus logger configured for JSON output:

```go
import "github.com/sirupsen/logrus"

func NewLogger() *logrus.Logger {
    l := logrus.New()
    l.SetFormatter(&logrus.JSONFormatter{})
    return l
}
```

### cobra

Even for non-CLI project types, selecting this add-on adds a `cmd/root.go` skeleton with a Cobra root command. This is useful for services that expose both an HTTP server and maintenance sub-commands in the same binary.

---

## AI

> Available on **all** project types. Adds a thin LLM client wrapper for ad-hoc AI calls without scaffolding a full AI Agent project.

### langchaingo

Adds `internal/ai/client.go` — a minimal LangChainGo client wrapper:

```go
import "github.com/tmc/langchaingo/llms/openai"

func NewLLMClient() (llms.Model, error) {
    return openai.New() // reads OPENAI_API_KEY from env
}
```

**After generation:**

```bash
go mod tidy   # pulls in github.com/tmc/langchaingo
```

Swap the provider import to use a different backend (e.g. `llms/googleai`, `llms/ollama`) with the same interface.

---

## Vector Store

> Primarily used with the **AI Agent** project type for Retrieval-Augmented Generation (RAG), but can be added to any project type.

### pgvector

Wires a [pgvector](https://github.com/pgvector/pgvector) store backed by PostgreSQL. Adds `internal/vectorstore/pgvector.go`:

```go
import (
    "github.com/tmc/langchaingo/vectorstores/pgvector"
    "github.com/jackc/pgx/v5/pgxpool"
)

func NewStore(ctx context.Context, pool *pgxpool.Pool, embedder embeddings.Embedder) (pgvector.Store, error) {
    return pgvector.New(ctx, pgvector.WithPool(pool), pgvector.WithEmbedder(embedder))
}
```

Requires a Postgres instance with the `pgvector` extension enabled:

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

**Best for:** Teams already running Postgres who want to add vector search without a separate infrastructure dependency.

### chromem

Adds `internal/vectorstore/chromem.go` using [chromem-go](https://github.com/philippgille/chromem-go) — an in-process, zero-dependency vector store.

```go
import "github.com/philippgille/chromem-go"

db := chromem.NewDB()
collection, err := db.GetOrCreateCollection("docs", nil, nil)
```

No external service required — the store lives in-process (optionally persisted to disk).

**Best for:** Local development, testing, small-scale production workloads where simplicity beats scalability.

### qdrant

Wires a [Qdrant](https://qdrant.tech) vector database via the official Go client. Adds `internal/vectorstore/qdrant.go`:

```go
import "github.com/tmc/langchaingo/vectorstores/qdrant"

store, err := qdrant.New(
    qdrant.WithURL(url.URL{Scheme: "http", Host: "localhost:6334"}),
    qdrant.WithCollectionName("my-collection"),
    qdrant.WithEmbedder(embedder),
)
```

Run Qdrant locally via Docker:

```bash
docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
```

**Best for:** Production-scale RAG systems that need filtering, payload storage, and a fully-featured vector database.

---

## Vector store comparison

| Store | Infrastructure | Scale | Persistence | Best for |
|-------|---------------|-------|-------------|----------|
| pgvector | PostgreSQL | Medium | Yes | Existing Postgres users |
| chromem  | In-process | Small | Optional | Local dev / testing |
| Qdrant   | Qdrant server | Large | Yes | Production RAG |

---

## Combining add-ons

All add-on files are independent of each other. Selecting `gorm` + `redis` + `zap` generates three separate files under `pkg/` — there is no conflict.

The main entrypoint (`main.go`) imports each package to ensure the code compiles immediately after `go mod tidy`.


## Selecting add-ons

In the web UI, add-ons appear in Section 04. You can combine add-ons freely across categories.

From the CLI:

```bash
goini new --addon cache=redis --addon database=gorm --addon other=zap
```

---

## Cache

### redis

Adds a `pkg/cache/redis.go` file with a configured `go-redis` client:

```go
import "github.com/redis/go-redis/v9"

func NewRedisClient(addr string) *redis.Client {
    return redis.NewClient(&redis.Options{Addr: addr})
}
```

**After generation:**

```bash
go mod tidy   # pulls in github.com/redis/go-redis/v9
```

Point the client at your Redis instance by passing the address (e.g. `localhost:6379`).

### memcached

Adds a `pkg/cache/memcached.go` file using `bradfitz/gomemcache`:

```go
import "github.com/bradfitz/gomemcache/memcache"

func NewMemcachedClient(addr ...string) *memcache.Client {
    return memcache.New(addr...)
}
```

---

## Database

### gorm

Adds `pkg/database/gorm.go` with a GORM connection helper and a sample model:

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

func NewDB(dsn string) (*gorm.DB, error) {
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

Switch the driver import (`sqlite`, `mysql`, `sqlserver`) to match your database.

### ent

Adds the Ent schema directory (`ent/schema/`) and a client initialiser. Run the Ent code-generator after customising your schema:

```bash
go run entgo.io/ent/cmd/ent generate ./ent/schema
```

---

## Other

### zap

Adds `pkg/logger/zap.go` with a production-ready Zap logger:

```go
import "go.uber.org/zap"

func NewLogger() (*zap.Logger, error) {
    return zap.NewProduction()
}
```

Use `zap.NewDevelopment()` locally for human-readable output.

### logrus

Adds `pkg/logger/logrus.go` with a Logrus logger configured for JSON output:

```go
import "github.com/sirupsen/logrus"

func NewLogger() *logrus.Logger {
    l := logrus.New()
    l.SetFormatter(&logrus.JSONFormatter{})
    return l
}
```

### cobra

Even for non-CLI project types, selecting this add-on adds a `cmd/root.go` skeleton with a Cobra root command. This is useful for services that expose both an HTTP server and maintenance sub-commands in the same binary.

---

## Combining add-ons

All add-on files are independent of each other. Selecting `gorm` + `redis` + `zap` generates three separate files under `pkg/` — there is no conflict.

The main entrypoint (`main.go`) imports each package to ensure the code compiles immediately after `go mod tidy`.
