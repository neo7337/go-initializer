# How to Use the Web UI

This guide walks through every section of the generator form.

## 1. Open the generator

Navigate to the app root. The generator form is the default view. Click **Explore** in the header to read documentation without leaving the app.

---

## 2. Select a Go Version

Pick the Go runtime version for your project. The latest stable version is shown first.

> **Tip:** Unless you have a specific reason to pin an older version, choose the highest available — it gives you access to the latest stdlib improvements and security patches.

---

## 3. Select a Project Type

| Option | Best for |
| --- | --- |
| Simple Project | Minimal single-binary service |
| Microservice | Service in a larger distributed system |
| API Server | Public-facing REST API |
| CLI Application | Terminal tools and automation scripts |

Choosing a project type unlocks the relevant framework options in the next section.

---

## 4. Select a Framework

The framework list updates automatically when you change the project type. Pick the HTTP framework or CLI library you want the scaffold to wire up.

See [Choosing a Framework](frameworks) for a detailed comparison.

---

## 5. Add-ons (optional)

Add-ons inject additional wiring into your project. They are grouped by category:

**Cache**
- `redis` — Redis client using `go-redis`
- `memcached` — Memcached client

**Database**
- `gorm` — GORM ORM with a SQLite/Postgres driver stub
- `ent` — Facebook Ent schema-first ORM

**Other**
- `zap` — Uber Zap structured logger
- `logrus` — Logrus structured logger
- `cobra` — Cobra CLI sub-command support (for non-CLI project types)

You can mix and match freely. If you select both a cache and a database, both are wired into the generated project.

---

## 6. Docker Support (optional)

Toggle **Generate Dockerfile** to include a production multi-stage `Dockerfile` in the project zip. See [Docker Setup](docker-setup) for details on what gets generated.

---

## 7. Project Details

| Field | Example | Notes |
| --- | --- | --- |
| Module Name | `github.com/acme/my-app` | Used in `go.mod` — follow Go module path conventions |
| Name | `my-app` | Used as the root directory name in the zip |
| Description | `Order management service` | Written into README and go files as a comment |

All three fields are required.

---

## 8. Generate

Click **Generate Project** (or press `⌘ Enter` on macOS / `Ctrl Enter` on other systems). A `project.zip` file downloads immediately.

```bash
# Unzip and get started
unzip project.zip
cd my-app
go mod tidy
go run .
```

---

## Keyboard shortcut

| Action | macOS | Windows / Linux |
| --- | --- | --- |
| Generate project | `⌘ ↵` | `Ctrl ↵` |
