package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prakoso-id/personal-backend/internal/config"
	"github.com/prakoso-id/personal-backend/internal/database"
	"github.com/prakoso-id/personal-backend/internal/middleware"
	"github.com/prakoso-id/personal-backend/internal/routes"
)

// @title           Personal Website API
// @version         1.0
// @description     Backend API for Personal Website.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect Database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate Database
	database.Migrate(db)

	// Setup Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())

	// Routes
	routes.RegisterRoutes(r, db, cfg)

	// Run Server
	log.Printf("Server running on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
