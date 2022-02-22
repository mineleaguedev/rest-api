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

	generalDB, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		os.Getenv("generaldb.username"),
		os.Getenv("generaldb.password"),
		os.Getenv("generaldb.host"),
		os.Getenv("generaldb.port"),
		os.Getenv("generaldb.dbname")),
	)
	if err != nil {
		log.Fatalf("Error connecting to general database: %s", err)
	}

	miniGamesDB, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		os.Getenv("minigamesdb.username"),
		os.Getenv("minigamesdb.password"),
		os.Getenv("minigamesdb.host"),
		os.Getenv("minigamesdb.port"),
		os.Getenv("minigamesdb.dbname")),
	)
	if err != nil {
		log.Fatalf("Error connecting to minigames database: %s", err)
	}

	controllers.Controller(generalDB, miniGamesDB)

	controllers.ConfigureCaptcha(os.Getenv("hcaptcha.site.key"), os.Getenv("hcaptcha.secret.key"))
	router.GET("/reg", controllers.RenderRegForm)
	router.GET("/auth", controllers.RenderAuthForm)
	router.POST("/reg", controllers.RegisterUser)
	router.POST("/auth", controllers.AuthUser)

	router.POST("/user", controllers.CreateUser)
	router.GET("/user/name/:name", controllers.GetUser)
	router.PUT("/user/exp", controllers.UpdateUserExp)
	router.PUT("/user/rank", controllers.UpdateUserRank)
	router.PUT("/user/playtime", controllers.UpdateUserPlaytime)
	router.PUT("/user/lastSeen", controllers.UpdateUserLastSeen)
	router.POST("/ban", controllers.BanUser)
	router.POST("/unban", controllers.UnbanUser)
	router.POST("/mute", controllers.MuteUser)
	router.POST("/unmute", controllers.UnmuteUser)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
