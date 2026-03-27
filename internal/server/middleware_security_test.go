package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newSecHeadersEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeadersMiddleware)
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestSecurityHeadersMiddleware_SetsAllHeaders(t *testing.T) {
	r := newSecHeadersEngine()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
	assert.Equal(t, "geolocation=(), microphone=(), camera=()", w.Header().Get("Permissions-Policy"))
	assert.Equal(t, "default-src 'none'", w.Header().Get("Content-Security-Policy"))
}

func TestSecurityHeadersMiddleware_NoHSTSHeader(t *testing.T) {
	// HSTS must not be set at the application layer — it is NPM's responsibility.
	r := newSecHeadersEngine()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Empty(t, w.Header().Get("Strict-Transport-Security"))
}

func TestSecurityHeadersMiddleware_HeadersPresentOnErrorResponses(t *testing.T) {
	// Headers should appear even when a downstream middleware aborts (e.g. 429).
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeadersMiddleware)
	r.GET("/test", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusTooManyRequests)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
}

func TestSecurityHeadersMiddleware_PresentOnHealthzViaFullRouter(t *testing.T) {
	// Verify headers flow through the real router on a plain GET.
	r := NewRouter(nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "default-src 'none'", w.Header().Get("Content-Security-Policy"))
}
