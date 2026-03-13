package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func MetaHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"supportedProjectTypes": SupportedProjectTypesLabelsMap,
		"supportedGoVersions":   SupportedGoVersionsMap,
		"supportedFrameworks":   SupportedFrameworksMap,
		"supportedAddons":       SupportedAddonsMap,
	})
}

func GenerateHandler(ctx *gin.Context) {
	var request CreateProjectRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
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
	log.Printf("[INFO] Received request: %+v", request)

	gen, ok := generatorRegistry[request.ProjectType]
	if !ok {
		log.Printf("[ERROR] Unsupported project type: %s", request.ProjectType)
		resp := ErrorResponseBody{
			StatusCode: http.StatusBadRequest,
			Message:    "Unsupported project type",
		}
		resp.GenerateResponse(ctx)
		return
	}

	buf, err := gen.Generate(request)
	if err != nil {
		log.Printf("[ERROR] Failed to generate project: %v", err)
		resp := ErrorResponseBody{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to generate project",
		}
		resp.GenerateResponse(ctx)
		return
	}
	resp := SuccessResponseBody{Data: buf.Bytes()}
	resp.GenerateResponse(ctx)
}
