package main

import (
	"log"

	"github.com/its-ayush-07/go-neux-server/initializers"
	"github.com/its-ayush-07/go-neux-server/models"
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
	initializers.DB.AutoMigrate(&models.Article{}, &models.User{})
}
