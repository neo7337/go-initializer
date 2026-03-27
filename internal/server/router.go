package server

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// trustedProxiesFromEnv reads the TRUSTED_PROXIES environment variable
// (comma-separated list of IP addresses or CIDRs, e.g. "172.18.0.0/16").
// Returns nil when absent, which callers treat as "no trusted proxies".
func trustedProxiesFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// allowedOriginsFromEnv reads the ALLOWED_ORIGINS environment variable
// (comma-separated list of origins). Returns nil when the variable is absent
// or blank, which signals wildcard mode to the caller.
func allowedOriginsFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func NewRouter(v *validator.Validate) *gin.Engine {
	service := gin.Default()

	// Configure trusted proxies so Gin reads X-Forwarded-For only from known
	// upstream proxies (e.g. the nginx proxy manager container on proxy-tier).
	// This prevents IP-spoofing via a crafted X-Forwarded-For header, which
	// would otherwise let an attacker circumvent the per-IP rate limit.
	trustedProxies := trustedProxiesFromEnv()
	if len(trustedProxies) > 0 {
		if err := service.SetTrustedProxies(trustedProxies); err != nil {
			log.Printf("[WARN] [SECURITY] Could not set trusted proxies (%v): %v — falling back to no trusted proxies", trustedProxies, err)
			_ = service.SetTrustedProxies(nil)
		}
	} else {
		// No TRUSTED_PROXIES configured: don't trust any proxy header. Gin will
		// use RemoteAddr directly. Log a notice so operators know.
		_ = service.SetTrustedProxies(nil)
		log.Println("[WARN] [SECURITY] TRUSTED_PROXIES is not configured — X-Forwarded-For headers will be ignored. " +
			"Set TRUSTED_PROXIES to the proxy-tier CIDR in production.")
	}

	// Security response headers — registered first so all responses including
	// 429 / 403 from downstream middleware carry the defensive headers.
	service.Use(SecurityHeadersMiddleware)

	// Rate limiting — applied before CORS so abusive callers are rejected early.
	rps, burst := rateLimitParamsFromEnv()
	service.Use(RateLimitMiddleware(newRateLimiterStore(rps, burst)))

	// Build CORS configuration from environment. If ALLOWED_ORIGINS is unset
	// the API defaults to the permissive wildcard, which is acceptable for a
	// fully-public, credential-free API but should be restricted in production.
	// AllowCredentials must not be used with a wildcard origin (browsers reject
	// such responses). Since this API serves file downloads with no cookies or
	// auth headers, credentials are not needed.
	allowedOrigins := allowedOriginsFromEnv()

	var corsConfig cors.Config
	if len(allowedOrigins) == 0 {
		log.Println("[WARN] [SECURITY] ALLOWED_ORIGINS is not configured — CORS wildcard (*) is active. " +
			"Set ALLOWED_ORIGINS to a comma-separated list of allowed origins in production.")
		corsConfig = cors.Config{
			AllowOrigins:  []string{"*"},
			AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:  []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "X-User-ID"},
			ExposeHeaders: []string{"Content-Length"},
			MaxAge:        12 * time.Hour,
		}
	} else {
		originSet := make(map[string]struct{}, len(allowedOrigins))
		for _, o := range allowedOrigins {
			originSet[o] = struct{}{}
		}
		// AllowOriginFunc returning false causes gin-contrib/cors to abort
		// preflight (OPTIONS) requests with 403 Forbidden for unlisted origins.
		corsConfig = cors.Config{
			AllowOriginFunc: func(origin string) bool {
				_, ok := originSet[origin]
				return ok
			},
			AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:  []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "X-User-ID"},
			ExposeHeaders: []string{"Content-Length"},
			MaxAge:        12 * time.Hour,
		}
	}

	service.Use(cors.New(corsConfig))

	service.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Server is up and running")
	})

	service.GET("/api/meta", MetaHandler)

	// MaxBodyBytesMiddleware is scoped to this route only — not the global chain —
	// because only /api/generate accepts a request body.
	service.POST("/api/generate", MaxBodyBytesMiddleware, GenerateHandler)

	return service
}
