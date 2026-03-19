# How to Self-Host Go Initializer

This guide covers deploying Go Initializer on your own infrastructure — a VPS, a Kubernetes cluster, or a local server — so your team can use it on an internal network.

## Architecture recap

Two services need to run:

| Service | Default port | Technology |
| --- | --- | --- |
| Backend API | `8181` | Go / Gin |
| Frontend | `8001` | Nginx (static files) |

---

## Docker Compose (simplest)

```bash
git clone https://github.com/neo7337/go-initializer.git
cd go-initializer

# Build and start both services
docker compose up --build -d

# View logs
docker compose logs -f

# Stop
docker compose down
```

---

## Behind a reverse proxy (Nginx)

Place an Nginx reverse proxy in front of both services to terminate TLS and serve everything on port 443.

```nginx
server {
    listen 443 ssl;
    server_name goini.example.com;

    ssl_certificate     /etc/ssl/certs/goini.crt;
    ssl_certificate_key /etc/ssl/private/goini.key;

    # Frontend (static files)
    location / {
        proxy_pass http://localhost:8001;
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8181;
    }

    location /healthz {
        proxy_pass http://localhost:8181;
    }
}
```

---

## Kubernetes

A minimal Kubernetes deployment for each service:

**backend-deployment.yaml**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-initializer-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: go-initializer-backend
  template:
    metadata:
      labels:
        app: go-initializer-backend
    spec:
      containers:
        - name: backend
          image: ghcr.io/neo7337/go-initializer-backend:latest
          ports:
            - containerPort: 8181
          env:
            - name: GIN_MODE
              value: release
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8181
            initialDelaySeconds: 5
            periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: go-initializer-backend
spec:
  selector:
    app: go-initializer-backend
  ports:
    - port: 8181
      targetPort: 8181
```

---

## Health checks

Both container images expose health endpoints for orchestrators:

```
GET /healthz   →  200 OK   (process is alive)
```

---

## Security recommendations for production

- Set `GIN_MODE=release` to disable verbose Gin debug output
- Place the backend behind a reverse proxy — do not expose it directly to the internet
- Restrict CORS origins in the backend's CORS config to your actual frontend domain
- Serve the frontend and backend over TLS (certificates from Let's Encrypt are free)
- Limit container capabilities with `securityContext.readOnlyRootFilesystem: true`
