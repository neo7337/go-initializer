# Troubleshooting

Solutions to common problems when running or using Go Initializer.

---

## The frontend shows "Error fetching metadata"

**Cause:** The frontend cannot reach the backend API.

**Fix:**

1. Make sure the backend is running:
   ```bash
   curl http://localhost:8182/healthz
   ```
2. Check that `VITE_API_URL` matches the backend address. If you started the backend on a non-default port, rebuild the frontend with the correct value:
   ```bash
   VITE_API_URL=http://localhost:9000 npm run build
   ```
3. If using Docker Compose, verify both containers started successfully:
   ```bash
   docker compose ps
   docker compose logs backend
   ```

---

## `go mod tidy` fails after unzipping the project

**Cause:** The Go proxy cannot reach the declared dependency.

**Fix:**

1. Check your internet connection and proxy settings (`GOPROXY`, `GONOSUMCHECK`).
2. If you are behind a corporate proxy, set:
   ```bash
   GOPROXY=https://goproxy.io,direct go mod tidy
   ```
3. Ensure the Go version installed locally matches or exceeds the version selected in the generator.

---

## The generated binary panics on startup

**Cause:** A required environment variable (e.g. database DSN, Redis address) is not set.

**Fix:** Check the generated `main.go` and `pkg/` files for any `os.Getenv` calls — they are the runtime dependencies. Set the required variables before running:

```bash
DATABASE_DSN="host=localhost user=app dbname=mydb" go run .
```

---

## Docker build fails: `no such file or directory`

**Cause:** The build context or Dockerfile path is wrong.

**Fix:** Run the build from the **repo root**, not from inside a subdirectory:

```bash
# Correct
docker build -f cmd/server/Dockerfile .

# Wrong (missing context)
cd cmd/server && docker build .
```

---

## `goini new` prompts hang or look garbled

**Cause:** The terminal does not support the interactive TUI (e.g. basic CI terminals, Emacs shell buffers).

**Fix:** Use flags to bypass all prompts:

```bash
goini new \
  --name myapp \
  --module github.com/acme/myapp \
  --type microservice \
  --framework gin \
  --go-version 1.25.0
```

---

## Browser downloads an empty or corrupt zip

**Cause:** A validation error was returned but the browser tried to interpret the JSON error body as a zip.

**Fix:** Open the browser DevTools → Network tab, repeat the action, and check the response body of the `/api/generate` request. Look for a JSON error message and fix the form field indicated in the `fields` key.

---

## Port already in use

```
listen tcp :8182: bind: address already in use
```

**Fix:** Find and kill the process using the port:

```bash
lsof -i :8182 | grep LISTEN
kill <PID>
```

Or start the server on a different port:

```bash
PORT=9000 go run ./cmd/server/...
```

---

## Still stuck?

- Open an issue on the [GitHub repository](https://github.com/neo7337/go-initializer/issues)
- Ask in the [Gophers Slack](https://gophers.slack.com) `#help` channel
