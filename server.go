package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/mineleaguedev/rest-api/controllers"
	"github.com/mineleaguedev/rest-api/general"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
	"os"
	"time"
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

	middleware := models.JWTMiddleware{
		Realm:            "mineleague jwt",
		SigningAlgorithm: "HS256",
		Key:              []byte(os.Getenv("jwt.secret.key")),
		AccessTokenTime:  1 * time.Minute,
		RefreshTokenTime: 30 * 24 * time.Hour,
		IdentityKey:      "username",
		TokenLookup:      "cookie:token",
		TokenHeadName:    "Bearer",

		SendCookie:     true,
		CookieMaxAge:   30 * 24 * time.Hour,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieDomain:   "localhost:8080",
		CookieName:     "token",
		CookieSameSite: http.SameSiteStrictMode,
	}

	controllers.Controller(generalDB, miniGamesDB)
	general.Setup(generalDB, middleware, os.Getenv("hcaptcha.site.key"), os.Getenv("hcaptcha.secret.key"))

	auth := router.Group("/auth")
	auth.GET("/reg", general.RenderRegForm)
	auth.GET("/auth", general.RenderAuthForm)
	auth.POST("/reg", general.RegHandler)
	auth.POST("/auth", general.AuthHandler)

	auth.GET("/refresh", general.RefreshHandler)
	auth.Use(general.MiddlewareFunc())
	{
		auth.GET("/logout", general.LogoutHandler)
	}

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
