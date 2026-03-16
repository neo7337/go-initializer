package generator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponseBody struct {
	Data []byte
}

func (r *SuccessResponseBody) GenerateResponse(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", "attachment; filename=project.zip")
	ctx.Data(http.StatusOK, "application/zip", r.Data)
}

type ErrorResponseBody struct {
	StatusCode int
	Message    string
}

func (r *ErrorResponseBody) GenerateResponse(ctx *gin.Context) {
	ctx.JSON(r.StatusCode, gin.H{"error": r.Message})
}
