package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func NewRouter(v *validator.Validate) *gin.Engine {
	service := gin.Default()

	// Configure CORS with more permissive settings for development
	service.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Common React dev server ports
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	service.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Server is up and running")
	})

	service.GET("/meta", MetaHandler)

	service.POST("/generate", GenerateHandler)

	return service
}
