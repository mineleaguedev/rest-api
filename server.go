package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	router := gin.Default()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	if _, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		os.Getenv("db.username"),
		os.Getenv("db.password"),
		os.Getenv("db.host"),
		os.Getenv("db.port"),
		os.Getenv("db.dbname")),
	); err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
