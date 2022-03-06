package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis.addr"),
		Password: os.Getenv("redis.password"),
		DB:       0,
	})
	_, err = redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Error connecting to redis: %s", err)
	}

	middleware := models.JWTMiddleware{
		Realm:            "mineleague jwt",
		AccessTokenKey:   []byte(os.Getenv("jwt.access.key")),
		RefreshTokenKey:  []byte(os.Getenv("jwt.refresh.key")),
		AccessTokenTime:  15 * time.Minute,
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
	general.Setup(generalDB, middleware, redisClient, os.Getenv("hcaptcha.site.key"), os.Getenv("hcaptcha.secret.key"))

	auth := router.Group("/auth")
	{
		auth.GET("/reg", general.RenderRegForm)
		auth.GET("/auth", general.RenderAuthForm)
		auth.POST("/reg", general.RegHandler)
		auth.POST("/auth", general.AuthHandler)
		auth.POST("/refresh", general.RefreshHandler)
		auth.POST("/logout", general.LogoutHandler)
	}

	router.Use(general.AuthMiddleware())
	{

	}

	api := router.Group("/api")
	{
		api.POST("/user", controllers.CreateUser)
		api.GET("/user/name/:name", controllers.GetUser)
		api.PUT("/user/exp", controllers.UpdateUserExp)
		api.PUT("/user/rank", controllers.UpdateUserRank)
		api.PUT("/user/playtime", controllers.UpdateUserPlaytime)
		api.PUT("/user/lastSeen", controllers.UpdateUserLastSeen)
		api.POST("/ban", controllers.BanUser)
		api.POST("/unban", controllers.UnbanUser)
		api.POST("/mute", controllers.MuteUser)
		api.POST("/unmute", controllers.UnmuteUser)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
