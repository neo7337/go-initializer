# Configuration Reference

All configuration is done through environment variables. No config files are required.

## Backend

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8182` | TCP port the HTTP server listens on |
| `GIN_MODE` | `debug` | Gin run mode. Set to `release` in production to silence debug output |

### Setting variables

**Shell:**
```bash
PORT=9000 GIN_MODE=release go run ./cmd/server/...
```

**Docker Compose:**
```yaml
services:
  backend:
    environment:
      - GIN_MODE=release
      - PORT=8181
```

**Kubernetes:**
```yaml
env:
  - name: GIN_MODE
    value: release
  - name: PORT
    value: "8181"
```

---

## Frontend

The frontend is configured **at build time** via Vite's environment variable mechanism. Variables must be prefixed with `VITE_`.

| Variable | Default | Description |
| --- | --- | --- |
| `VITE_API_URL` | `http://localhost:8182` | Full origin of the backend API, used for all fetch calls |

### Setting for local development

Create a `.env.local` file in the `frontend/` directory:

```
VITE_API_URL=http://localhost:8182
```

Vite automatically loads `.env.local` (it is gitignored by default).

### Setting for a production build

```bash
VITE_API_URL=https://api.goini.example.com npm run build
```

Or via Docker build arg:

```dockerfile
ARG VITE_API_URL=https://api.goini.example.com
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build
```

---

## CORS

The backend allows all origins (`*`) by default. This is suitable for development and for a fully public tool, but if you self-host and only want your own frontend to call the API, restrict the origins:

```go
// internal/server/router.go
cors.New(cors.Config{
    AllowOrigins: []string{"https://goini.example.com"},
    ...
})
```

---

## Timeouts

The frontend applies a **15-second request timeout** to all API calls (configurable in `frontend/src/service.ts`):

```ts
const DEFAULT_TIMEOUT_MS = 15_000;
```

For very large projects or slow servers, increase this value in your fork.
