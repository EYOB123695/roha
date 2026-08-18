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

	_ "github.com/EYOB123695/roha/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

// @title Roha Social API
// @version 1.0
// @description Backend API for Roha social platform.
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Create a Gin router
	r := gin.Default()

	// Initializing Infrastructure / Repositories
	userRepo := repository.NewUserRepository(initializers.DB)
	postRepo := repository.NewPostRepository(initializers.DB)
	commentRepo := repository.NewCommentRepository(initializers.DB)

	// Initializing Use Cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	postUseCase := usecase.NewPostUseCase(postRepo)
	commentUseCase := usecase.NewCommentUseCase(commentRepo, postRepo)

	// Initializing Handlers (HTTP Adapters)
	userHandler := httpDelivery.NewUserHandler(userUseCase)
	postHandler := httpDelivery.NewPostHandler(postUseCase)
	commentHandler := httpDelivery.NewCommentHandler(commentUseCase)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "roha API is running in Clean Architecture",
		})
	})

	// loader.io verification routes
	loaderTokenHandler := func(c *gin.Context) {
		c.String(200, "loaderio-2243a8a4089d49a8edc7a2eec362c2b8")
	}
	r.GET("/loaderio-2243a8a4089d49a8edc7a2eec362c2b8", loaderTokenHandler)
	r.GET("/loaderio-2243a8a4089d49a8edc7a2eec362c2b8.txt", loaderTokenHandler)
	r.GET("/loaderio-2243a8a4089d49a8edc7a2eec362c2b8.html", loaderTokenHandler)
	r.GET("/loaderio-2243a8a4089d49a8edc7a2eec362c2b8/", loaderTokenHandler)

	// Public Routes
	r.POST("/signup", userHandler.Signup)
	r.POST("/login", userHandler.Login)
	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:id", postHandler.GetPost)
	r.GET("/posts/:id/comments", commentHandler.GetCommentsByPostID)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})
    
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
		protected.POST("/posts/:id/comments", commentHandler.AddComment)
		protected.DELETE("/comments/:id", commentHandler.DeleteComment)
	    protected.GET("/feed", postHandler.GetFeed) 
		protected.POST("/posts/:id/like", postHandler.LikePost) 
		// <-- Add this route
		protected.DELETE("/posts/:id/like", postHandler.UnlikePost)
		protected.POST("/posts/:id/bookmark", postHandler.BookmarkPost)
		protected.DELETE("/posts/:id/bookmark", postHandler.UnbookmarkPost)
		protected.GET("/bookmarks", postHandler.GetBookmarkedPosts)
		protected.POST("/posts/:id/activity", postHandler.TrackActivity)
		protected.GET("/posts/recommendations", postHandler.GetRecommendations)

	}

	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
