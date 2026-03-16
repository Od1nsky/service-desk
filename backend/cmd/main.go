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
	ticketRepo := repository.NewTicketRepo(database)
	commentRepo := repository.NewCommentRepo(database)
	categoryRepo := repository.NewCategoryRepo(database)
	priorityRepo := repository.NewPriorityRepo(database)
	statsRepo := repository.NewStatsRepo(database)

	// Handlers
	authH := handlers.NewAuthHandler(userRepo, cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	ticketH := handlers.NewTicketHandler(ticketRepo)
	commentH := handlers.NewCommentHandler(commentRepo)
	categoryH := handlers.NewCategoryHandler(categoryRepo)
	priorityH := handlers.NewPriorityHandler(priorityRepo)
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
		// Tickets
		protected.GET("/tickets", ticketH.List)
		protected.POST("/tickets", middleware.RequireRole("employee"), ticketH.Create)
		protected.GET("/tickets/:id", ticketH.GetByID)
		protected.PATCH("/tickets/:id/status", middleware.RequireRole("support", "admin"), ticketH.UpdateStatus)
		protected.PATCH("/tickets/:id/assign", middleware.RequireRole("support"), ticketH.Assign)
		protected.PATCH("/tickets/:id/rate", middleware.RequireRole("employee"), ticketH.Rate)

		// Comments
		protected.GET("/tickets/:id/comments", commentH.List)
		protected.POST("/tickets/:id/comments", commentH.Create)

		// Categories
		protected.GET("/categories", categoryH.List)
		protected.POST("/categories", middleware.RequireRole("admin"), categoryH.Create)
		protected.PUT("/categories/:id", middleware.RequireRole("admin"), categoryH.Update)
		protected.DELETE("/categories/:id", middleware.RequireRole("admin"), categoryH.Delete)

		// Priorities
		protected.GET("/priorities", priorityH.List)

		// Users (admin)
		protected.GET("/users", middleware.RequireRole("admin"), userH.List)
		protected.PATCH("/users/:id/role", middleware.RequireRole("admin"), userH.UpdateRole)
		protected.PATCH("/users/:id/block", middleware.RequireRole("admin"), userH.Block)

		// Stats (admin + support)
		protected.GET("/stats/dashboard", middleware.RequireRole("admin", "support"), statsH.Dashboard)
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
