package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/mineleaguedev/rest-api/controllers"
	"log"
	"os"
)

func main() {
	router := gin.Default()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		os.Getenv("db.username"),
		os.Getenv("db.password"),
		os.Getenv("db.host"),
		os.Getenv("db.port"),
		os.Getenv("db.dbname")),
	)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	controllers.Controller(db)

	router.POST("/user", controllers.CreateUser)
	router.GET("/user/name/:name", controllers.GetUser)
	router.PUT("/user/exp", controllers.UpdateUserExp)
	router.PUT("/user/rank", controllers.UpdateUserRank)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
