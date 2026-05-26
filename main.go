package main

import (
	"log"
	"os"

	"github.com/EYOB123695/roha/controllers"
	"github.com/EYOB123695/roha/initializers"
	"github.com/EYOB123695/roha/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "roha API is running",
			"endpoints": gin.H{
				"create_post": "POST /post",
			},
		})
	})
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/posts", controllers.GetPosts)
	r.GET("/posts/:id", controllers.GetPost)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth)
	{
		protected.GET("/validate", controllers.Validate)
		protected.POST("/posts", controllers.PostsCreate)
		protected.PUT("/posts/:id", controllers.UpdatePost)
		protected.DELETE("/posts/:id", controllers.DeletePost)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server on port 8080 (default) or the env port
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
