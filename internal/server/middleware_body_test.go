package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// newBodyLimitEngine wires a minimal engine with MaxBodyBytesMiddleware and a
// handler that reads the full body via ShouldBindJSON (mirroring GenerateHandler).
func newBodyLimitEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/generate", MaxBodyBytesMiddleware, func(ctx *gin.Context) {
		var body map[string]interface{}
		if err := ctx.ShouldBindJSON(&body); err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "too large"})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.Status(http.StatusOK)
	})
	return r
}

func TestMaxBodyBytesMiddleware_AcceptsBodyUnderLimit(t *testing.T) {
	r := newBodyLimitEngine()

	body := `{"key":"` + strings.Repeat("a", 100) + `"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMaxBodyBytesMiddleware_Rejects413WhenBodyExceedsLimit(t *testing.T) {
	r := newBodyLimitEngine()

	// Build a JSON body well over 64 KB.
	bigValue := strings.Repeat("x", 70_000)
	body := `{"key":"` + bigValue + `"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestMaxBodyBytesMiddleware_ExactLimitIsAccepted(t *testing.T) {
	// A JSON body at exactly (or just under) the limit should be accepted.
	// Use limit - 8 for the value so the full JSON object fits within 64 KB.
	r := newBodyLimitEngine()

	padding := strings.Repeat("a", int(maxRequestBodyBytes)-8)
	body := `{"k":"` + padding + `"}`
	assert.LessOrEqual(t, int64(len(body)), maxRequestBodyBytes, "test body should be at or under the limit")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestMaxBodyBytesMiddleware_NotAppliedToOtherRoutes(t *testing.T) {
	// GET /healthz has no body limit — a large body is simply ignored.
	r := NewRouter(nil)

	bigBody := bytes.Repeat([]byte("x"), 70_000)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", bytes.NewReader(bigBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
