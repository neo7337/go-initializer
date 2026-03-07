# Microservices Best Practices

Go Initializer's microservices template follows a layout that scales from a handful of services to dozens.

## Recommended folder structure

```
myapp/
├── services/
│   ├── auth/
│   │   ├── main.go
│   │   ├── handler.go
│   │   └── Dockerfile
│   └── orders/
│       ├── main.go
│       ├── handler.go
│       └── Dockerfile
├── pkg/
│   ├── models/
│   └── middleware/
├── docker-compose.yml
└── go.work
```

`pkg/` holds code shared across services. Each service is an independent binary with its own `main.go`.

## Use Go Workspaces

Go 1.18+ workspaces (`go.work`) let you develop multiple modules in a single repo without `replace` directives:

```bash
go work init ./services/auth ./services/orders
```

## Keep services focused

Each service should own a single bounded context. Resist the urge to share database models across service boundaries — communicate via APIs or events instead.

## Health & readiness endpoints

Every service should expose:

- `GET /healthz` — liveness probe (returns 200 if the process is alive)
- `GET /readyz` — readiness probe (returns 200 when dependencies are ready)

## Structured logging

Use **Zap** or **slog** (stdlib, Go 1.21+) for structured JSON logs that work well with log aggregators:

```go
logger, _ := zap.NewProduction()
logger.Info("server started", zap.Int("port", 8080))
```
