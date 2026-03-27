package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- allowedOriginsFromEnv tests ---

func TestAllowedOriginsFromEnv_EmptyReturnsNil(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "")
	assert.Nil(t, allowedOriginsFromEnv())
}

func TestAllowedOriginsFromEnv_SingleOrigin(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://example.com")
	got := allowedOriginsFromEnv()
	assert.Equal(t, []string{"https://example.com"}, got)
}

func TestAllowedOriginsFromEnv_MultipleOriginsWithSpaces(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", " https://a.com , https://b.com ")
	got := allowedOriginsFromEnv()
	assert.Equal(t, []string{"https://a.com", "https://b.com"}, got)
}

// --- CORS behaviour tests ---

func TestNewRouter_PreflightFromAllowedOriginSucceeds(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://allowed.com")
	r := NewRouter(validator.New())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/generate", nil)
	req.Header.Set("Origin", "https://allowed.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusForbidden, w.Code)
	assert.Equal(t, "https://allowed.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestNewRouter_PreflightFromDisallowedOriginReturns403(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://allowed.com")
	r := NewRouter(validator.New())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/generate", nil)
	req.Header.Set("Origin", "https://attacker.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestNewRouter_WildcardModeAllowsAnyOrigin(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "")
	r := NewRouter(validator.New())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/generate", nil)
	req.Header.Set("Origin", "https://anyone.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusForbidden, w.Code)
}

func TestNewRouter_ReturnsEngineAndHealthRoute(t *testing.T) {
	r := NewRouter(validator.New())
	require.NotNil(t, r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Server is up and running")
}

func TestNewRouter_UnknownRouteReturns404(t *testing.T) {
	r := NewRouter(validator.New())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNewRouter_RequestWithOriginHeaderStillSucceeds(t *testing.T) {
	r := NewRouter(validator.New())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "https://example.com")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewRouter_NilValidatorStillBuildsRouter(t *testing.T) {
	// NewRouter currently does not use the validator parameter directly.
	r := NewRouter(nil)
	require.NotNil(t, r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
