package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	router := gin.Default()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
