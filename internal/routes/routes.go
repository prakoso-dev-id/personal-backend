package routes

import (	
	"github.com/gin-gonic/gin"
	"github.com/prakoso-id/personal-backend/internal/config"
	"github.com/prakoso-id/personal-backend/internal/middleware"
	"github.com/prakoso-id/personal-backend/internal/modules/auth"
	"github.com/prakoso-id/personal-backend/internal/modules/contact"
	"github.com/prakoso-id/personal-backend/internal/modules/images"
	"github.com/prakoso-id/personal-backend/internal/modules/experiences"
	"github.com/prakoso-id/personal-backend/internal/modules/posts"
	"github.com/prakoso-id/personal-backend/internal/modules/profiles"
	"github.com/prakoso-id/personal-backend/internal/modules/projects"
	"github.com/prakoso-id/personal-backend/internal/modules/skills"
    
    // Swagger
    "github.com/prakoso-id/personal-backend/docs"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize base URL for image paths
	images.SetBaseURL(cfg.Server.BaseURL)

	// Swagger Info
	docs.SwaggerInfo.BasePath = "/api"

	// Repositories
	authRepo := auth.NewRepository(db)
	imageRepo := images.NewRepository(db)
	profileRepo := profiles.NewRepository(db)
	postRepo := posts.NewRepository(db)
	projectRepo := projects.NewRepository(db)
	experienceRepo := experiences.NewRepository(db)

	// Services
	authService := auth.NewService(authRepo, cfg)
	imageService := images.NewService(imageRepo)
	profileService := profiles.NewService(profileRepo)
	postService := posts.NewService(postRepo, imageRepo)
	projectService := projects.NewService(projectRepo, imageRepo)
	experienceService := experiences.NewService(experienceRepo)

	// Handlers
	authHandler := auth.NewHandler(authService)
	imageHandler := images.NewHandler(imageService)
	profileHandler := profiles.NewHandler(profileService)
	postHandler := posts.NewHandler(postService)
	projectHandler := projects.NewHandler(projectService)
	skillHandler := skills.NewHandler(db)
	contactHandler := contact.NewHandler(db)
	experienceHandler := experiences.NewHandler(experienceService, profileService)

	api := r.Group("/api")
	{
		// Public Routes
		public := api.Group("/public")
		{
			public.GET("/profile", profileHandler.GetProfile)
			public.GET("/skills", skillHandler.GetAll)
			public.GET("/posts", postHandler.GetPublicPosts)
			public.GET("/posts/:slug", postHandler.GetPublicPostBySlug)
			public.GET("/projects", projectHandler.GetPublicProjects)
			public.GET("/projects/:id", projectHandler.GetPublicProjectByID)
			public.GET("/experiences", experienceHandler.GetPublicExperiences)
			public.POST("/contact", contactHandler.CreateMessage)

			// Swagger
			public.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}

		// Admin Routes (Protected)
		admin := api.Group("/admin")
		admin.POST("/login", authHandler.Login)
		
		protected := admin.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Profile (Admin)
			protected.GET("/profile", profileHandler.GetProfile)
			protected.PUT("/profile", profileHandler.UpdateProfile)

			// Auth Updates
			protected.PUT("/update-email", authHandler.UpdateEmail)
			protected.PUT("/update-password", authHandler.UpdatePassword)

			// Posts (Admin)
			protected.GET("/posts", postHandler.GetAdminPosts)
			protected.POST("/posts", postHandler.CreatePost)
			protected.PUT("/posts/:id", postHandler.UpdatePost)
			protected.DELETE("/posts/:id", postHandler.DeletePost)

			// Projects (Admin)
			protected.GET("/projects", projectHandler.GetAdminProjects)
			protected.POST("/projects", projectHandler.CreateProject)
			protected.PUT("/projects/:id", projectHandler.UpdateProject)
			protected.DELETE("/projects/:id", projectHandler.Delete)

			// Skills (Admin)
			protected.GET("/skills", skillHandler.GetAll)
			protected.POST("/skills", skillHandler.Create)
			protected.PUT("/skills/:id", skillHandler.Update)
			protected.DELETE("/skills/:id", skillHandler.Delete)

			// Images
			protected.POST("/images/upload", imageHandler.Upload)
			protected.DELETE("/images/:id", imageHandler.Delete)
			
			// Contact Messages (Read)
			protected.GET("/messages", contactHandler.GetAllMessages)

			// Experiences (Admin)
			protected.GET("/experiences", experienceHandler.GetAdminExperiences)
			protected.POST("/experiences", experienceHandler.CreateExperience)
			protected.PUT("/experiences/:id", experienceHandler.UpdateExperience)
			protected.DELETE("/experiences/:id", experienceHandler.DeleteExperience)
		}
	}
	
	// Static file serving for images
	// Map /media to storage folder
	// In production this might be handled by Nginx
	r.Static("/media", "./storage")
}
