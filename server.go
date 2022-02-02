package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
