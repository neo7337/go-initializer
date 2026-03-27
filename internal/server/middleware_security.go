package server

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware sets defensive HTTP response headers on every
// response served by the backend API.
//
// Header rationale:
//   - X-Content-Type-Options: nosniff       — prevents MIME-type sniffing
//   - X-Frame-Options: DENY                 — disallows embedding in frames
//   - Referrer-Policy                       — limits referrer leakage
//   - Permissions-Policy                    — disables browser features not needed by this API
//   - Content-Security-Policy: default-src 'none'
//     The backend serves only JSON and binary (zip) responses — no HTML, no
//     scripts, no resources. A blanket deny is the tightest possible policy.
//
// HSTS is intentionally omitted: in the neolabs-infra deployment, TLS is
// terminated by nginx proxy manager (NPM) which issues Let's Encrypt certificates.
// HSTS should be configured in the NPM proxy host Advanced tab:
//
//	add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
func SecurityHeadersMiddleware(ctx *gin.Context) {
	ctx.Header("X-Content-Type-Options", "nosniff")
	ctx.Header("X-Frame-Options", "DENY")
	ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
	ctx.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	ctx.Header("Content-Security-Policy", "default-src 'none'")
	ctx.Next()
}
