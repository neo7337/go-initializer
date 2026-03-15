package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func NewRouter(v *validator.Validate) *gin.Engine {
	service := gin.Default()

	// Configure CORS. AllowCredentials must not be used with a wildcard origin
	// (browsers will reject the response). Since this API serves file downloads
	// with no cookies or auth headers, credentials are not needed.
	service.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "X-User-ID"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	service.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Server is up and running")
	})

	service.GET("/api/meta", MetaHandler)

	service.POST("/api/generate", GenerateHandler)

	return service
}
