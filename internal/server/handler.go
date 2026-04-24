package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/neo7337/go-initializer/internal/generator"
)

func MetaHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"supportedProjectTypes": generator.SupportedProjectTypesLabelsMap,
		"supportedGoVersions":   generator.SupportedGoVersionsMap,
		"supportedFrameworks":   generator.SupportedFrameworksMap,
		"supportedAddons":       generator.SupportedAddonsMap,
	})
}

func GenerateHandler(ctx *gin.Context) {
	var request generator.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		// Detect an oversized body before reporting a generic bad-request error.
		// *http.MaxBytesError is set by MaxBodyBytesMiddleware via MaxBytesReader.
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body exceeds the 64 KB limit"})
			return
		}
		log.Printf("[ERROR] Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := validate.Struct(request); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make(map[string]string, len(ve))
			for _, fe := range ve {
				fields[fe.Field()] = fe.Tag()
			}
			log.Printf("[WARN] Validation failed: %v", fields)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "fields": fields})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	for category, values := range request.Addons {
		if len(values) > 1 {
			log.Printf("[WARN] Multiple addons for category %q: %v", category, values)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("addon category %q allows at most one selection", category),
			})
			return
		}
	}

	log.Printf("[INFO] Received request: projectType=%s name=%s moduleName=%s goVersion=%s framework=%s",
		request.ProjectType, request.Name, request.ModuleName, request.GoVersion, request.Framework)

	gen, ok := generator.GeneratorRegistry[request.ProjectType]
	if !ok {
		log.Printf("[ERROR] Unsupported project type: %s", request.ProjectType)
		resp := generator.ErrorResponseBody{
			StatusCode: http.StatusBadRequest,
			Message:    "Unsupported project type",
		}
		resp.GenerateResponse(ctx)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	buf, err := gen.Generate(timeoutCtx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("[ERROR] Generate timed out for project type: %s", request.ProjectType)
			ctx.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out generating project"})
			return
		}
		var ve *generator.ErrValidation
		if errors.As(err, &ve) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": ve.Error()})
			return
		}
		log.Printf("[ERROR] Failed to generate project: %v", err)
		resp := generator.ErrorResponseBody{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to generate project",
		}
		resp.GenerateResponse(ctx)
		return
	}
	resp := generator.SuccessResponseBody{Data: buf.Bytes()}
	resp.GenerateResponse(ctx)
}
