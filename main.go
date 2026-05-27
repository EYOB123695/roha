package main

import (
	"log"
	"os"

	httpDelivery "github.com/EYOB123695/roha/delivery/http"
	"github.com/EYOB123695/roha/initializers"
	"github.com/EYOB123695/roha/middleware"
	"github.com/EYOB123695/roha/repository"
	"github.com/EYOB123695/roha/usecase"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	// Create a Gin router
	r := gin.Default()

	// Initializing Infrastructure / Repositories
	userRepo := repository.NewUserRepository(initializers.DB)
	postRepo := repository.NewPostRepository(initializers.DB)

	// Initializing Use Cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	postUseCase := usecase.NewPostUseCase(postRepo)

	// Initializing Handlers (HTTP Adapters)
	userHandler := httpDelivery.NewUserHandler(userUseCase)
	postHandler := httpDelivery.NewPostHandler(postUseCase)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "roha API is running in Clean Architecture",
		})
	})

	// Public Routes
	r.POST("/signup", userHandler.Signup)
	r.POST("/login", userHandler.Login)
	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:id", postHandler.GetPost)

	// Protected Routes (Uses injected requireAuth middleware)
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth(userRepo))
	{
		protected.GET("/validate", userHandler.Validate)
		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
		protected.GET("/users/:id", userHandler.GetUserProfile)
		protected.POST("/users/:id/follow", userHandler.FollowUser)
		protected.POST("/users/:id/unfollow", userHandler.UnfollowUser)
		protected.GET("/users/:id/followers", userHandler.GetFollowers)
		protected.GET("/users/:id/following", userHandler.GetFollowing)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
