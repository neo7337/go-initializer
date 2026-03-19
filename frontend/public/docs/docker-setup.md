# Docker Setup

Go Initializer ships a ready-made `docker-compose.yml` and per-service `Dockerfile`s so you can get a full stack running with a single command.

## Quick start

```bash
git clone https://github.com/neo7337/go-initializer.git
cd go-initializer

docker compose up --build
```

- Frontend: `http://localhost:8001`
- Backend API: `http://localhost:8181`

To run in the background:

```bash
docker compose up --build -d
```

To stop and remove containers:

```bash
docker compose down
```

---

## Services

### backend

| Property | Value |
| --- | --- |
| Build context | repo root |
| Dockerfile | `cmd/server/Dockerfile` |
| Exposed port | `8181` |
| Environment | `GIN_MODE=release` |

The backend image uses a multi-stage build: a `golang:1.24` builder stage compiles the binary, and a minimal `gcr.io/distroless/static` image is used for the final artifact — no shell, no package manager, minimal attack surface.

### frontend

| Property | Value |
| --- | --- |
| Build context | `./frontend` |
| Dockerfile | `frontend/Dockerfile` |
| Exposed port | `8001` |
| Served by | Nginx |

The frontend image runs `npm run build` in a Node.js builder stage and then serves the static output via Nginx.

---

## Customising the Compose file

### Change ports

Edit `docker-compose.yml`:

```yaml
services:
  backend:
    ports:
      - "9000:8181"   # host:container
  frontend:
    ports:
      - "3000:8001"
```

### Point the frontend at a different backend

Pass the backend URL as a build argument if you're rebuilding the frontend image:

```yaml
services:
  frontend:
    build:
      args:
        VITE_API_URL: http://my-backend.example.com
```

Or set `VITE_API_URL` before running `npm run build` locally.

### Add a reverse proxy

To terminate TLS and expose everything on port 443, place an Nginx or Caddy reverse proxy in front of the `frontend` service. The backend does **not** need to be publicly exposed when a proxy handles routing.

---

## Generated project Docker support

When you enable **Docker Support** in the generator, the downloaded zip includes a production-grade multi-stage `Dockerfile` for your new project:

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

Build and run it from your generated project directory:

```bash
docker build -t my-app .
docker run -p 8080:8080 my-app
```
