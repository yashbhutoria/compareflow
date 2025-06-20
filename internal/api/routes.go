package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/compareflow/compareflow/internal/api/handlers"
	"github.com/compareflow/compareflow/internal/api/middleware"
	"github.com/compareflow/compareflow/internal/config"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// CORS middleware
	router.Use(middleware.CORS(cfg.AllowedOrigins))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "compareflow",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Authentication routes (no auth required)
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	v1.POST("/auth/login", authHandler.Login)
	v1.POST("/auth/register", authHandler.Register)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))

	// User routes
	protected.GET("/auth/me", authHandler.Me)

	// Connection routes
	connectionHandler := handlers.NewConnectionHandler(db)
	protected.GET("/connections", connectionHandler.List)
	protected.GET("/connections/:id", connectionHandler.Get)
	protected.POST("/connections", connectionHandler.Create)
	protected.PUT("/connections/:id", connectionHandler.Update)
	protected.DELETE("/connections/:id", connectionHandler.Delete)
	protected.POST("/connections/:id/test", connectionHandler.Test)
	protected.GET("/connections/:id/tables", connectionHandler.GetTables)
	protected.GET("/connections/:id/tables/:table/columns", connectionHandler.GetColumns)

	// Validation routes
	validationHandler := handlers.NewValidationHandler(db)
	protected.GET("/validations", validationHandler.List)
	protected.GET("/validations/:id", validationHandler.Get)
	protected.POST("/validations", validationHandler.Create)
	protected.PUT("/validations/:id", validationHandler.Update)
	protected.DELETE("/validations/:id", validationHandler.Delete)
	protected.POST("/validations/:id/run", validationHandler.Run)
	protected.GET("/validations/:id/status", validationHandler.Status)
}