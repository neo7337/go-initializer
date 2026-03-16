package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
