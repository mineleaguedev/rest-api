package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/kataras/hcaptcha"
	"github.com/mineleaguedev/rest-api/handlers"
	"github.com/mineleaguedev/rest-api/models"
	"github.com/mineleaguedev/rest-api/services"
	"github.com/nitishm/go-rejson/v4"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
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

	// setup jwt
	middleware := models.JWTMiddleware{
		Realm:              "mineleague jwt",
		AccessTokenKey:     []byte(os.Getenv("jwt.access.key")),
		RefreshTokenKey:    []byte(os.Getenv("jwt.refresh.key")),
		AccessTokenTime:    15 * time.Minute,
		RefreshTokenTime:   30 * 24 * time.Hour,
		RegTokenTime:       30 * 60,
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

	// setup aws
	awsRegion := aws.String(os.Getenv("aws.region"))

	// setup email
	sess, err := session.NewSession(&aws.Config{
		Region: awsRegion,
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("aws.ses.access.key.id"),
			os.Getenv("aws.ses.secret.access.key"),
			"",
		),
	})
	if err != nil {
		log.Fatalf("Error connecting to email smtp: %s", err)
	}
	emailClient := ses.New(sess)

	// setup skin
	sess, err = session.NewSession(&aws.Config{
		Region: awsRegion,
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("aws.s3.skins.access.key.id"),
			os.Getenv("aws.s3.skins.secret.access.key"),
			"",
		),
	})
	if err != nil {
		log.Fatalf("Error connecting to skins s3: %s", err)
	}
	skinUploader := s3manager.NewUploader(sess)
	skinDeleter := s3.New(sess)

	// setup cloak
	sess, err = session.NewSession(&aws.Config{
		Region: awsRegion,
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("aws.s3.cloaks.access.key.id"),
			os.Getenv("aws.s3.cloaks.secret.access.key"),
			"",
		),
	})
	if err != nil {
		log.Fatalf("Error connecting to cloaks s3: %s", err)
	}
	cloakUploader := s3manager.NewUploader(sess)
	cloakDeleter := s3.New(sess)

	service := services.NewService(
		middleware,
		models.RedisConfig{
			Client:      redisClient,
			JsonHandler: redisJsonHandler,
			Ctx:         context.TODO(),
		}, models.EmailConfig{
			RegFrom:            viper.GetString("email.reg.from"),
			RegSubject:         viper.GetString("email.reg.subject"),
			RegHtmlBody:        viper.GetString("email.reg.htmlBody"),
			PassResetFrom:      viper.GetString("email.passReset.from"),
			PassResetSubject:   viper.GetString("email.passReset.subject"),
			PassResetHtmlBody:  viper.GetString("email.passReset.htmlBody"),
			NewPassFrom:        viper.GetString("email.newPass.from"),
			NewPassSubject:     viper.GetString("email.newPass.subject"),
			NewPassHtmlBody:    viper.GetString("email.newPass.htmlBody"),
			ChangePassFrom:     viper.GetString("email.changePass.from"),
			ChangePassSubject:  viper.GetString("email.changePass.subject"),
			ChangePassHtmlBody: viper.GetString("email.changePass.htmlBody"),
			Client:             emailClient,
		}, models.CaptchaConfig{
			SiteKey:         os.Getenv("hcaptcha.site.key"),
			Client:          hcaptcha.New(os.Getenv("hcaptcha.secret.key")),
			RegForm:         template.Must(template.ParseFiles("./forms/reg_form.html")),
			AuthForm:        template.Must(template.ParseFiles("./forms/auth_form.html")),
			PassResetForm:   template.Must(template.ParseFiles("./forms/pass_reset_form.html")),
			ChangePassForm:  template.Must(template.ParseFiles("./forms/change_pass_form.html")),
			ChangeSkinForm:  template.Must(template.ParseFiles("./forms/change_skin_form.html")),
			DeleteSkinForm:  template.Must(template.ParseFiles("./forms/delete_skin_form.html")),
			ChangeCloakForm: template.Must(template.ParseFiles("./forms/change_cloak_form.html")),
			DeleteCloakForm: template.Must(template.ParseFiles("./forms/delete_cloak_form.html")),
		}, models.SkinConfig{
			SkinBucket:    aws.String(os.Getenv("aws.s3.skins.bucket.name")),
			SkinUploader:  skinUploader,
			SkinDeleter:   skinDeleter,
			CloakBucket:   aws.String(os.Getenv("aws.s3.cloaks.bucket.name")),
			CloakUploader: cloakUploader,
			CloakDeleter:  cloakDeleter,
		})
	handler := handlers.NewHandler(service, middleware, generalDB, miniGamesDB)

	router := handler.InitRoutes()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
