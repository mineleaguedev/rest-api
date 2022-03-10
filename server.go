package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/mineleaguedev/rest-api/controllers"
	"github.com/mineleaguedev/rest-api/general"
	"github.com/mineleaguedev/rest-api/models"
	"github.com/nitishm/go-rejson/v4"
	"github.com/spf13/viper"
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

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	middleware := models.JWTMiddleware{
		Realm:              "mineleague jwt",
		AccessTokenKey:     []byte(os.Getenv("jwt.access.key")),
		RefreshTokenKey:    []byte(os.Getenv("jwt.refresh.key")),
		AccessTokenTime:    15 * time.Minute,
		RefreshTokenTime:   30 * 24 * time.Hour,
		RegTokenTime:       30 * 1000,
		PassResetTokenTime: 30 * time.Minute,
		IdentityKey:        "username",
		TokenLookup:        "cookie:token",
		TokenHeadName:      "Bearer",

		SendCookie:     true,
		CookieMaxAge:   30 * 24 * time.Hour,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieDomain:   "localhost:8080",
		CookieName:     "token",
		CookieSameSite: http.SameSiteStrictMode,
	}

	// setup email
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("aws.ses.region")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("aws.ses.access.key.id"),
			os.Getenv("aws.ses.secret.access.key"),
			"",
		),
	})
	emailClient := ses.New(sess)
	general.SetupEmail(emailClient)

	// setup redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis.addr"),
		Password: os.Getenv("redis.password"),
		DB:       0,
	})
	_, err = redisClient.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatalf("Error connecting to redis: %s", err)
	}
	redisJsonHandler := rejson.NewReJSONHandler()
	redisJsonHandler.SetGoRedisClient(redisClient)
	general.SetupRedis(redisClient, redisJsonHandler)

	// setup hcaptcha
	hcaptchaSiteKey := os.Getenv("hcaptcha.site.key")
	hcaptchaSecretKey := os.Getenv("hcaptcha.secret.key")
	general.SetupCaptcha(hcaptchaSiteKey, hcaptchaSecretKey)

	controllers.Controller(generalDB, miniGamesDB)
	general.Setup(generalDB, middleware)

	auth := router.Group("/auth")
	{
		auth.GET("/reg", general.RenderRegForm)
		auth.POST("/reg", general.RegHandler)
		auth.GET("/reg/confirm/:token", general.ConfirmRegHandler)
		auth.GET("/auth", general.RenderAuthForm)
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
