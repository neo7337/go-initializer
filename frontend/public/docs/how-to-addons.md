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

## Combining add-ons

All add-on files are independent of each other. Selecting `gorm` + `redis` + `zap` generates three separate files under `pkg/` — there is no conflict.

The main entrypoint (`main.go`) imports each package to ensure the code compiles immediately after `go mod tidy`.
