# Installation

This page covers every way to run Go Initializer — using the hosted web UI, self-hosting the server locally, running with Docker, or installing the `goini` CLI.

## Prerequisites

| Requirement | Minimum version | Notes |
| --- | --- | --- |
| Go | 1.24 | Required for local server build and CLI install |
| Node.js | 18 | Required only if you want to build the frontend yourself |
| Docker | 24 | Required only for the Docker Compose workflow |
| Git | any | Required to clone the repository |

---

## Option 1 — Use the hosted UI

Visit the live deployment and start generating immediately. No installation required.

---

## Option 2 — Run locally (from source)

```bash
# 1. Clone the repository
git clone https://github.com/neo7337/go-initializer.git
cd go-initializer

# 2. Install Go dependencies
go mod tidy

# 3. Start the backend (listens on :8182)
make dev-server
# or directly:
go run ./cmd/server/...

# 4. In a second terminal, start the frontend (listens on :5173)
cd frontend
npm install
npm start
```

Open `http://localhost:5173` in your browser.

---

## Option 3 — Docker Compose (recommended for self-hosting)

The repository ships a `docker-compose.yml` that builds and runs both the backend and frontend containers:

```bash
git clone https://github.com/neo7337/go-initializer.git
cd go-initializer

docker compose up --build
```

| Service | Port | URL |
| --- | --- | --- |
| Backend | 8181 | `http://localhost:8181` |
| Frontend | 8001 | `http://localhost:8001` |

To stop:

```bash
docker compose down
```

---

## Option 4 — Install the goini CLI

The `goini` binary lets you scaffold projects directly from your terminal without opening a browser.

### From source

```bash
git clone https://github.com/neo7337/go-initializer.git
cd go-initializer

# Build the CLI binary to ./bin/goini
make build-cli

# Move it to a directory in your PATH
sudo mv bin/goini /usr/local/bin/goini
```

### Verify

```bash
goini version
# goini dev   (or the release tag)
```

---

## Building the frontend for production

If you want to serve the frontend as a static bundle (e.g. from an Nginx container):

```bash
cd frontend
npm install
npm run build
```

The output is written to `frontend/build/`. The included `frontend/Dockerfile` performs this build automatically.

---

## Environment variables

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8182` | Port the backend server listens on |
| `GIN_MODE` | `debug` | Set to `release` in production to disable Gin's debug output |
| `VITE_API_URL` | `http://localhost:8182` | Backend origin used by the frontend at build time |
