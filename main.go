package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/its-ayush-07/go-neux-server/controllers"
	"github.com/its-ayush-07/go-neux-server/initializers"
	"github.com/its-ayush-07/go-neux-server/middleware"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	initializers.ConnectToDB()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = false
	config.AllowOrigins = append(config.AllowOrigins, "http://localhost:3000")
	config.AllowOrigins = append(config.AllowOrigins, os.Getenv("CLIENT_ORIGIN"))
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Routes
	r.GET("/articles", controllers.RecentArticles)
	r.GET("/articles/:id", controllers.ArticleByID)
	r.GET("/search", controllers.SearchArticles)
	r.GET("/user/:id", controllers.UserByID)
	r.GET("/user/articles/:id", controllers.ArticleByAuthor)
	r.POST("/user/signup", controllers.SignUp)
	r.POST("/user/login", controllers.Login)
	r.Use(middleware.Authentication())
	r.POST("/user/create", controllers.CreateArticle)
	r.PATCH("/articles/:id/like", controllers.UpdateLikes)
	r.DELETE("/articles/:id", controllers.DeleteArticleByID)
	r.POST("/user/logout", controllers.Logout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default to port 8000 if not specified
	}

	// Start the server
	addr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
