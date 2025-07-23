package main

import (
	"github.com/chris910512/travel-chat/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if envErr := godotenv.Load(); envErr != nil {
		log.Println("No .env file found")
	}

	db, dbErr := database.NewPostgresDB()
	if dbErr != nil {
		log.Fatal("Failed to connect to database:", dbErr)
	}

	if migErr := database.RunMigrations(db); migErr != nil {
		log.Fatal("Failed to migrate database:", migErr)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	runErr := r.Run(":8080")
	if runErr != nil {
		return
	}
}
