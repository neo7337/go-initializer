# REST API Reference

The Go Initializer backend exposes three endpoints. All endpoints are served over HTTP; there is no authentication.

## Base URL

```
http://localhost:8182
```

Change the port via the `PORT` environment variable.

---

## GET /healthz

Liveness probe. Returns `200 OK` when the process is running.

**Response**

```
200 OK
"Server is up and running"
```

---

## GET /api/meta

Returns every valid value for the generator form fields. Call this endpoint once on app load to populate drop-downs.

**Response schema**

```json
{
  "supportedGoVersions": {
    "1.25.0": true,
    "1.24.6": true,
    "1.23.12": true
  },
  "supportedProjectTypes": {
    "microservice": "Microservice",
    "simple-project": "Simple Project",
    "cli-app": "CLI Application",
    "api-server": "API Server"
  },
  "supportedFrameworks": {
    "microservice": { "gin": true, "echo": true, "fiber": true, "golly": true, "gokit": true },
    "api-server":   { "gin": true, "echo": true, "fiber": true, "chi": true, "golly": true },
    "cli-app":      { "cobra": true, "urfave": true, "kingpin": true, "golly": true },
    "simple-project": { "golly": true }
  },
  "supportedAddons": {
    "cache":    { "redis": true, "memcached": true },
    "database": { "gorm": true, "ent": true },
    "other":    { "zap": true, "logrus": true, "cobra": true }
  }
}
```

---

## POST /api/generate

Generates a Go project and returns it as a zip archive.

**Request headers**

```
Content-Type: application/json
```

**Request body**

```json
{
  "projectType":   "microservice",
  "goVersion":     "1.25.0",
  "framework":     "gin",
  "moduleName":    "github.com/acme/order-service",
  "name":          "order-service",
  "description":   "Order management service",
  "selectedAddons": {
    "cache":    ["redis"],
    "database": ["gorm"],
    "other":    ["zap"]
  },
  "dockerSupport": true
}
```

**Field reference**

| Field | Type | Required | Valid values |
| --- | --- | --- | --- |
| `projectType` | string | yes | `microservice` · `simple-project` · `cli-app` · `api-server` |
| `goVersion` | string | yes | `1.25.0` · `1.24.6` · `1.23.12` |
| `framework` | string | yes | depends on `projectType` — see `/api/meta` |
| `moduleName` | string | yes | any valid Go module path |
| `name` | string | yes | used as root directory name in the zip |
| `description` | string | no | free text |
| `selectedAddons` | object | no | map of category → array of addon values |
| `dockerSupport` | boolean | no | `true` to include a `Dockerfile` |

**Success response**

```
200 OK
Content-Type: application/zip
Content-Disposition: attachment; filename="project.zip"

<binary zip data>
```

**Error responses**

| Status | Body | Reason |
| --- | --- | --- |
| `400` | `{"error": "Invalid request body"}` | Malformed JSON |
| `400` | `{"error": "validation failed", "fields": {...}}` | Required field missing or invalid value |
| `400` | `{"error": "Unsupported project type"}` | `projectType` not in the registry |
| `500` | `{"error": "Failed to generate project"}` | Internal generation error |

---

## Example: cURL

```bash
curl -X POST http://localhost:8182/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "projectType": "api-server",
    "goVersion": "1.25.0",
    "framework": "gin",
    "moduleName": "github.com/acme/my-api",
    "name": "my-api",
    "description": "My new API",
    "dockerSupport": true
  }' \
  --output my-api.zip
```
