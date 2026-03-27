package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// newTestEngine wires a minimal Gin engine with the supplied middleware and
// a single GET /test handler that always returns 200.
func newTestEngine(mw gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw)
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

// --- rateLimitParamsFromEnv ---

func TestRateLimitParamsFromEnv_Defaults(t *testing.T) {
	t.Setenv("RATE_LIMIT_RPS", "")
	t.Setenv("RATE_LIMIT_BURST", "")
	rps, burst := rateLimitParamsFromEnv()
	assert.Equal(t, 5.0, rps)
	assert.Equal(t, 10, burst)
}

func TestRateLimitParamsFromEnv_CustomValues(t *testing.T) {
	t.Setenv("RATE_LIMIT_RPS", "20")
	t.Setenv("RATE_LIMIT_BURST", "50")
	rps, burst := rateLimitParamsFromEnv()
	assert.Equal(t, 20.0, rps)
	assert.Equal(t, 50, burst)
}

func TestRateLimitParamsFromEnv_InvalidValuesUseDefaults(t *testing.T) {
	t.Setenv("RATE_LIMIT_RPS", "notanumber")
	t.Setenv("RATE_LIMIT_BURST", "-5")
	rps, burst := rateLimitParamsFromEnv()
	assert.Equal(t, 5.0, rps)
	assert.Equal(t, 10, burst)
}

// --- rateLimiterStore ---

func TestGetLimiter_ReturnsSameInstanceForSameIP(t *testing.T) {
	store := newRateLimiterStore(5, 10)
	l1 := store.getLimiter("1.2.3.4")
	l2 := store.getLimiter("1.2.3.4")
	assert.Same(t, l1, l2)
}

func TestGetLimiter_ReturnsDifferentInstanceForDifferentIP(t *testing.T) {
	store := newRateLimiterStore(5, 10)
	l1 := store.getLimiter("1.2.3.4")
	l2 := store.getLimiter("5.6.7.8")
	assert.NotSame(t, l1, l2)
}

// --- RateLimitMiddleware ---

func TestRateLimitMiddleware_AllowsRequestsUnderBurst(t *testing.T) {
	// Large burst so no request will be rejected during the test.
	store := newRateLimiterStore(100, 100)
	r := newTestEngine(RateLimitMiddleware(store))

	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "request %d should pass", i+1)
	}
}

func TestRateLimitMiddleware_Blocks429WhenBurstExhausted(t *testing.T) {
	// burst=1, rps very low — second request for same IP must be blocked.
	store := newRateLimiterStore(0.001, 1)
	r := newTestEngine(RateLimitMiddleware(store))

	allowed := 0
	blocked := 0
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:4567"
		r.ServeHTTP(w, req)
		if w.Code == http.StatusOK {
			allowed++
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
			blocked++
		}
	}
	assert.Equal(t, 1, allowed, "exactly burst=1 request should be allowed")
	assert.Equal(t, 4, blocked, "remaining requests should be blocked")
}

func TestRateLimitMiddleware_RetryAfterHeaderPresent(t *testing.T) {
	store := newRateLimiterStore(0.001, 1)
	r := newTestEngine(RateLimitMiddleware(store))

	// First request drains the burst.
	r.ServeHTTP(httptest.NewRecorder(), newReqFromIP("10.1.1.1"))

	// Second request must be blocked with Retry-After.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, newReqFromIP("10.1.1.1"))

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.NotEmpty(t, w.Header().Get("Retry-After"))
}

func TestRateLimitMiddleware_DifferentIPsHaveIndependentLimits(t *testing.T) {
	// burst=1 means each IP gets exactly one free request.
	store := newRateLimiterStore(0.001, 1)
	r := newTestEngine(RateLimitMiddleware(store))

	for _, ip := range []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, newReqFromIP(ip))
		assert.Equal(t, http.StatusOK, w.Code, "first request from %s should be allowed", ip)
	}
}

// --- trustedProxiesFromEnv (in router.go) ---

func TestTrustedProxiesFromEnv_EmptyReturnsNil(t *testing.T) {
	t.Setenv("TRUSTED_PROXIES", "")
	assert.Nil(t, trustedProxiesFromEnv())
}

func TestTrustedProxiesFromEnv_SingleCIDR(t *testing.T) {
	t.Setenv("TRUSTED_PROXIES", "172.18.0.0/16")
	got := trustedProxiesFromEnv()
	assert.Equal(t, []string{"172.18.0.0/16"}, got)
}

func TestTrustedProxiesFromEnv_MultipleWithSpaces(t *testing.T) {
	t.Setenv("TRUSTED_PROXIES", " 10.0.0.1 , 10.0.0.2 ")
	got := trustedProxiesFromEnv()
	assert.Equal(t, []string{"10.0.0.1", "10.0.0.2"}, got)
}

// helpers

func newReqFromIP(ip string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = ip + ":9999"
	return req
}
