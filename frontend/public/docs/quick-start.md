# Quick Start

Get a Go project generated and running in under two minutes.

## 1. Choose a project type

Open the generator and pick between **Simple Project** or **Microservices**.

- **Simple Project** — ideal for a single HTTP service or CLI tool
- **Microservices** — generates a multi-service workspace with a shared `pkg/` directory

## 2. Fill in the details

| Field | Description |
|-------|-------------|
| Module Name | Your Go module path, e.g. `github.com/you/myapp` |
| Project Name | The name of the root directory |
| Go Version | Select the Go runtime version |
| Framework | HTTP framework (Gin, Echo, Fiber, …) |

## 3. Pick add-ons (optional)

Enable any combination of:

- **Databases** — GORM or Ent ORM wiring
- **Cache** — Redis or Memcached client setup
- **Logging** — Zap structured logger
- **CLI** — Cobra command scaffolding
- **Docker** — `Dockerfile` and `docker-compose.yml`

## 4. Generate & download

Click **Generate Project** — a `.zip` archive downloads immediately. Unzip and run:

```bash
cd your-project
go mod tidy
go run .
```

Your server is live on `localhost:8080`.
