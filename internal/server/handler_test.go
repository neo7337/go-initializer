package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// newTestRouter returns a fully wired router in test mode.
func newTestRouter() *gin.Engine {
	return NewRouter(validator.New())
}

// ─── /healthz ────────────────────────────────────────────────────────────────

func TestHealthz_Returns200(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── GET /api/meta ───────────────────────────────────────────────────────────

func TestMetaHandler_Returns200WithExpectedKeys(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/meta", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body, "supportedProjectTypes")
	assert.Contains(t, body, "supportedGoVersions")
	assert.Contains(t, body, "supportedFrameworks")
	assert.Contains(t, body, "supportedAddons")
}

func TestMetaHandler_SupportedProjectTypesNotEmpty(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/meta", nil)
	r.ServeHTTP(w, req)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	types, ok := body["supportedProjectTypes"].(map[string]interface{})
	require.True(t, ok, "supportedProjectTypes should be a map")
	assert.NotEmpty(t, types)
}

// ─── POST /api/generate ──────────────────────────────────────────────────────

func TestGenerateHandler_InvalidJSON(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate",
		bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body, "error")
}

func TestGenerateHandler_EmptyBody_ValidationFails(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate",
		bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body, "error")
	// Validation error includes the failing fields
	assert.Contains(t, body, "fields")
}

func TestGenerateHandler_MissingRequiredFields(t *testing.T) {
	r := newTestRouter()

	payload := map[string]interface{}{
		"projectType": "microservice",
		// missing: goVersion, framework, moduleName, name
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "validation failed", resp["error"])
	fields, ok := resp["fields"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, fields, "GoVersion")
}

func TestGenerateHandler_InvalidProjectType(t *testing.T) {
	r := newTestRouter()

	payload := map[string]interface{}{
		"projectType": "bad-type",
		"goVersion":   "1.24.6",
		"framework":   "gin",
		"moduleName":  "github.com/acme/x",
		"name":        "x",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGenerateHandler_InvalidGoVersion(t *testing.T) {
	r := newTestRouter()

	payload := map[string]interface{}{
		"projectType": "microservice",
		"goVersion":   "1.99.0", // not in oneof
		"framework":   "gin",
		"moduleName":  "github.com/acme/x",
		"name":        "x",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGenerateHandler_APIServerNotInRegistry_Returns400(t *testing.T) {
	// api-server passes the validator oneof check but is not in GeneratorRegistry
	r := newTestRouter()

	payload := map[string]interface{}{
		"projectType": "api-server",
		"goVersion":   "1.24.6",
		"framework":   "gin",
		"moduleName":  "github.com/acme/x",
		"name":        "x",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp, "error")
}

func TestGenerateHandler_SimpleProject_ReturnsZip(t *testing.T) {
	r := newTestRouter()

	payload := map[string]interface{}{
		"projectType": "simple-project",
		"goVersion":   "1.24.6",
		"framework":   "golly",
		"moduleName":  "github.com/acme/myproj",
		"name":        "myproj",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Body.Bytes())
}

func TestGenerateHandler_Microservice_ReturnsZip(t *testing.T) {
	r := newTestRouter()

	payload := map[string]interface{}{
		"projectType": "microservice",
		"goVersion":   "1.24.6",
		"framework":   "gin",
		"moduleName":  "github.com/acme/svc",
		"name":        "svc",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Body.Bytes())
}

func TestGenerateHandler_CLIApp_ReturnsZip(t *testing.T) {
	r := newTestRouter()

	payload := map[string]interface{}{
		"projectType": "cli-app",
		"goVersion":   "1.24.6",
		"framework":   "cobra",
		"moduleName":  "github.com/acme/tool",
		"name":        "tool",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Body.Bytes())
}

func TestGenerateHandler_SimpleProject_WithDocker_ReturnsZip(t *testing.T) {
	r := newTestRouter()

	payload := map[string]interface{}{
		"projectType":   "simple-project",
		"goVersion":     "1.24.6",
		"framework":     "golly",
		"moduleName":    "github.com/acme/myproj",
		"name":          "myproj",
		"dockerSupport": true,
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
}

func TestGenerateHandler_DispositionHeader(t *testing.T) {
	r := newTestRouter()

	payload := map[string]interface{}{
		"projectType": "simple-project",
		"goVersion":   "1.24.6",
		"framework":   "golly",
		"moduleName":  "github.com/acme/myproj",
		"name":        "myproj",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	disposition := w.Header().Get("Content-Disposition")
	assert.Contains(t, disposition, "attachment")
	assert.Contains(t, disposition, "project.zip")
}

func TestGenerateHandler_MethodNotAllowed_Get(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/generate", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
