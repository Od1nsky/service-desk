package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/servicedesk/backend/internal/config"
	"github.com/servicedesk/backend/internal/db"
	"github.com/servicedesk/backend/internal/handlers"
	"github.com/servicedesk/backend/internal/middleware"
	"github.com/servicedesk/backend/internal/repository"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer database.Close()

	// Repositories
	userRepo := repository.NewUserRepo(database)
	gradeRepo := repository.NewGradeRepo(database)
	commentRepo := repository.NewCommentRepo(database)
	disciplineRepo := repository.NewDisciplineRepo(database)
	gradeTypeRepo := repository.NewGradeTypeRepo(database)
	statsRepo := repository.NewStatsRepo(database)

	// Handlers
	authH := handlers.NewAuthHandler(userRepo, cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	gradeH := handlers.NewGradeHandler(gradeRepo)
	commentH := handlers.NewCommentHandler(commentRepo)
	disciplineH := handlers.NewDisciplineHandler(disciplineRepo)
	gradeTypeH := handlers.NewGradeTypeHandler(gradeTypeRepo)
	userH := handlers.NewUserHandler(userRepo)
	statsH := handlers.NewStatsHandler(statsRepo)

	r := gin.Default()

	// CORS
	r.Use(corsMiddleware(cfg.CORSOrigins))

	api := r.Group("/api/v1")

	// ---- Public auth routes ----
	auth := api.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", authH.Logout)
	}

	// ---- Authenticated routes ----
	protected := api.Group("")
	protected.Use(middleware.Auth(cfg.JWTSecret))
	{
		// Grades
		protected.GET("/grades", gradeH.List)
		protected.POST("/grades", middleware.RequireRole("student"), gradeH.Create)
		protected.GET("/grades/:id", gradeH.GetByID)
		protected.PATCH("/grades/:id/status", middleware.RequireRole("teacher", "admin"), gradeH.UpdateStatus)
		protected.PATCH("/grades/:id/assign", middleware.RequireRole("teacher"), gradeH.Assign)
		protected.PATCH("/grades/:id/acknowledge", middleware.RequireRole("teacher", "admin"), gradeH.Acknowledge)

		// Comments
		protected.GET("/grades/:id/comments", commentH.List)
		protected.POST("/grades/:id/comments", commentH.Create)

		// Disciplines
		protected.GET("/disciplines", disciplineH.List)
		protected.POST("/disciplines", middleware.RequireRole("admin"), disciplineH.Create)
		protected.PUT("/disciplines/:id", middleware.RequireRole("admin"), disciplineH.Update)
		protected.DELETE("/disciplines/:id", middleware.RequireRole("admin"), disciplineH.Delete)

		// Grade types
		protected.GET("/grade-types", gradeTypeH.List)

		// Users (admin)
		protected.GET("/users", middleware.RequireRole("admin"), userH.List)
		protected.PATCH("/users/:id/role", middleware.RequireRole("admin"), userH.UpdateRole)
		protected.PATCH("/users/:id/block", middleware.RequireRole("admin"), userH.Block)

		// Stats (admin + teacher)
		protected.GET("/stats/dashboard", middleware.RequireRole("admin", "teacher"), statsH.Dashboard)
	}

	log.Printf("starting server on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}

// corsMiddleware sets CORS headers for every request and handles OPTIONS preflight.
func corsMiddleware(origins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
