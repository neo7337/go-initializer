package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
	log.Printf("[INFO] Received request: %+v", request)
	switch request.ProjectType {
	case "microservice":
		buf, err := GenerateMicroservice(request)
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
	case "simple-project":
		buf, err := GenerateSimpleProject(request)
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
	default:
		log.Printf("[ERROR] Unsupported project type: %s", request.ProjectType)
		resp := ErrorResponseBody{
			StatusCode: http.StatusBadRequest,
			Message:    "Unsupported project type",
		}
		resp.GenerateResponse(ctx)
		return
	}
}
