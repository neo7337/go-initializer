package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// maxRequestBodyBytes is the hard cap on the JSON request body for
// /api/generate (64 KB). A well-formed CreateProjectRequest is under 2 KB;
// 64 KB gives plenty of headroom while preventing memory exhaustion from
// maliciously large payloads.
const maxRequestBodyBytes int64 = 1 << 16 // 64 KB

// MaxBodyBytesMiddleware wraps the incoming request body with
// http.MaxBytesReader so that reads beyond maxRequestBodyBytes return an
// *http.MaxBytesError. The handler is responsible for detecting that error
// and returning 413 Request Entity Too Large.
//
// Register this middleware only on the routes that accept a request body
// (e.g. POST /api/generate) — not on the global chain.
func MaxBodyBytesMiddleware(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxRequestBodyBytes)
	ctx.Next()
}
