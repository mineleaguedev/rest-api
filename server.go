package main

import (
	"database/sql"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/mineleaguedev/rest-api/controllers"
	"github.com/mineleaguedev/rest-api/systems"
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

	controllers.Controller(generalDB, miniGamesDB)
	systems.System(generalDB)

	systems.ConfigureCaptcha(os.Getenv("hcaptcha.site.key"), os.Getenv("hcaptcha.secret.key"))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "mineleague jwt",
		Key:         []byte(os.Getenv("jwt.secret.key")),
		Timeout:     15 * time.Minute,
		MaxRefresh:  30 * 24 * time.Hour,
		IdentityKey: "username",

		PayloadFunc:     systems.Payload,
		IdentityHandler: systems.Identify,
		Authenticator:   systems.Authenticate,
		Authorizator:    systems.Authorize,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieDomain:   "localhost:8080",
		CookieName:     "token",
		TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteDefaultMode,

		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatalf("Error initializing jwt: %s", err.Error())
	}
	if errInit := authMiddleware.MiddlewareInit(); errInit != nil {
		log.Fatalf("Error initializing systems middleware: %s", err.Error())
	}

	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/reg", systems.RegisterUser)

	auth := router.Group("/auth")
	auth.GET("/reg", systems.RenderRegForm)
	auth.GET("/auth", systems.RenderAuthForm)
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", systems.HelloHandler)
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
