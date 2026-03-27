package server

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter pairs a token-bucket limiter with a last-seen timestamp so the
// cleanup goroutine can evict stale entries.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// rateLimiterStore holds per-IP limiters and the shared rate parameters.
type rateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	rps      rate.Limit
	burst    int
}

func newRateLimiterStore(rps float64, burst int) *rateLimiterStore {
	s := &rateLimiterStore{
		limiters: make(map[string]*ipLimiter),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
	go s.cleanup()
	return s
}

// getLimiter returns (and lazily creates) the per-IP limiter.
func (s *rateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.limiters[ip]
	if !ok {
		entry = &ipLimiter{limiter: rate.NewLimiter(s.rps, s.burst)}
		s.limiters[ip] = entry
	}
	entry.lastSeen = time.Now()
	return entry.limiter
}

// cleanup removes limiters that have been idle for > 10 minutes to prevent
// unbounded memory growth when the server is exposed to many unique IPs.
func (s *rateLimiterStore) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		s.mu.Lock()
		for ip, entry := range s.limiters {
			if time.Since(entry.lastSeen) > 10*time.Minute {
				delete(s.limiters, ip)
			}
		}
		s.mu.Unlock()
	}
}

// rateLimitParamsFromEnv reads RATE_LIMIT_RPS and RATE_LIMIT_BURST.
// Defaults: rps = 5.0, burst = 10.
func rateLimitParamsFromEnv() (rps float64, burst int) {
	rps = 5.0
	burst = 10

	if v := os.Getenv("RATE_LIMIT_RPS"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			rps = f
		}
	}
	if v := os.Getenv("RATE_LIMIT_BURST"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			burst = i
		}
	}
	return rps, burst
}

// RateLimitMiddleware returns a per-IP token-bucket rate limiting middleware.
// The IP is obtained via ctx.ClientIP(), which respects the trusted-proxy
// configuration set on the Gin engine (T-SEC2 requirement).
// When a request exceeds the limit it is rejected with 429 Too Many Requests
// and a Retry-After header indicating the minimum seconds to wait.
// The [SECURITY] log tag allows operators to grep for rate events.
func RateLimitMiddleware(store *rateLimiterStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		lim := store.getLimiter(ip)

		if !lim.Allow() {
			// Retry-After: ceiling of (1 / rps) seconds.
			retryAfter := int(math.Ceil(1.0 / float64(store.rps)))
			if retryAfter < 1 {
				retryAfter = 1
			}
			log.Printf("[SECURITY] Rate limit exceeded for IP %s", ip)
			ctx.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded — please retry later",
			})
			return
		}

		ctx.Next()
	}
}
