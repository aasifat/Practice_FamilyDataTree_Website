package main

import (
	"fmt"
	"log"

	"family-tree-api/config"
	"family-tree-api/internal/database"
	"family-tree-api/internal/handlers"
	"family-tree-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.Init(config.AppConfig.GetDSN()); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(middleware.CORSMiddleware())

	// Public routes
	public := router.Group("/api/v1")
	{
		// Auth routes
		public.POST("/auth/signup", handlers.SignUp)
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/forgot-password", handlers.ForgotPassword)
		public.POST("/auth/reset-password", handlers.ResetPassword)
		public.POST("/auth/request-otp", handlers.RequestOTP)
		public.POST("/auth/verify-otp", handlers.VerifyOTPAndLogin)
		public.POST("/auth/request-otp-signup", handlers.RequestOTPForSignup)
		public.POST("/auth/verify-otp-signup", handlers.VerifyOTPAndSignup)

		// Contact routes
		public.POST("/contact", handlers.CreateMessage)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthRequired())
	{
		// Auth routes
		protected.GET("/auth/profile", handlers.GetProfile)

		// Family tree routes
		protected.POST("/trees", handlers.CreateFamilyTree)
		protected.GET("/trees", handlers.GetUserTrees)

		// Tree members routes (specific routes before wildcard)
		protected.GET("/trees/:id/members/search", handlers.SearchPeople)
		protected.GET("/trees/:id/members", handlers.GetTreeMembers)
		protected.POST("/trees/:id/members", handlers.CreatePerson)

		// Family tree detail routes (after more specific tree routes)
		protected.GET("/trees/:id", handlers.GetFamilyTree)
		protected.PUT("/trees/:id", handlers.UpdateFamilyTree)
		protected.DELETE("/trees/:id", handlers.DeleteFamilyTree)

		// Person routes
		protected.GET("/members/:personId", handlers.GetPerson)
		protected.PUT("/members/:personId", handlers.UpdatePerson)
		protected.DELETE("/members/:personId", handlers.DeletePerson)
		protected.GET("/members/:personId/children", handlers.GetChildren)
	}

	// Admin routes
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		// User management
		admin.GET("/users", handlers.GetAllUsers)
		admin.DELETE("/users/:id", handlers.DeleteUserByAdmin)

		// Messages
		admin.GET("/messages", handlers.GetAllMessages)
	}

	// Start server
	address := fmt.Sprintf(":%d", config.AppConfig.ServerPort)
	log.Printf("Server running on %s", address)
	router.Run(address)
}
