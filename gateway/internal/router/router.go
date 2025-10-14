package router

import (
	"gateway/internal/clients"
	"gateway/internal/handler"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(r *gin.Engine, grpcClients *clients.Clients, redis *redis.Client) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	api := r.Group("/api")
	{
		// Public routes
		public := api.Group("")
		{
			// Auth
			authHandler := handler.NewAuthHandler(grpcClients.Auth)
			public.POST("/auth/register", authHandler.Register)
			public.POST("/auth/login", authHandler.Login)
			public.POST("/auth/refresh", authHandler.RefreshToken)

			// News (public)
			newsHandler := handler.NewNewsHandler(grpcClients.News)
			public.GET("/news", newsHandler.ListNews)
			public.GET("/news/:slug", newsHandler.GetNewsBySlug)
			public.GET("/categories", newsHandler.ListCategories)
			public.GET("/tags", newsHandler.ListTags)

			// SEO
			seoHandler := handler.NewSEOHandler(grpcClients.SEO)
			public.GET("/sitemap.xml", seoHandler.GetSitemap)
			public.GET("/robots.txt", seoHandler.GetRobots)
		}

		// Protected routes (JWT required)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(grpcClients.Auth))
		{
			authHandler := handler.NewAuthHandler(grpcClients.Auth)
			protected.POST("/auth/logout", authHandler.Logout)

			// News management
			newsHandler := handler.NewNewsHandler(grpcClients.News)
			protected.POST("/news", middleware.RequireRole("editor", "admin"), newsHandler.CreateNews)
			protected.PUT("/news/:id", middleware.RequireRole("editor", "admin"), newsHandler.UpdateNews)
			protected.DELETE("/news/:id", middleware.RequireRole("admin"), newsHandler.DeleteNews)
			protected.POST("/news/:id/publish", middleware.RequireRole("editor", "admin"), newsHandler.PublishNews)

			// Categories management
			protected.POST("/categories", middleware.RequireRole("admin"), newsHandler.CreateCategory)

			// Tags management
			protected.POST("/tags", middleware.RequireRole("admin"), newsHandler.CreateTag)

			// Media
			mediaHandler := handler.NewMediaHandler(grpcClients.Media)
			protected.POST("/media/upload", mediaHandler.UploadFile)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(grpcClients.Auth))
		admin.Use(middleware.RequireRole("admin", "moderator"))
		{
			adminHandler := handler.NewAdminHandler(grpcClients.Admin)
			admin.GET("/users", adminHandler.GetUsers)
			admin.GET("/statistics", adminHandler.GetStatistics)
			admin.POST("/moderate/:id", adminHandler.ModerateContent)
		}
	}
}
