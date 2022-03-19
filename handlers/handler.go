package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/controllers"
	"github.com/mineleaguedev/rest-api/handlers/auth"
	"github.com/mineleaguedev/rest-api/handlers/cabinet"
	"github.com/mineleaguedev/rest-api/models"
	"github.com/mineleaguedev/rest-api/services"
)

type Handler struct {
	cabinet  *cabinet.Handler
	auth     *auth.Handler
	services *services.Service
}

func NewHandler(services *services.Service, middleware models.JWTMiddleware, generalDB *sql.DB) *Handler {
	return &Handler{
		cabinet:  cabinet.NewHandler(services, generalDB),
		auth:     auth.NewHandler(services, middleware, generalDB),
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/reg", h.auth.RegHandler)
		authGroup.GET("/reg/confirm/:token", h.auth.ConfirmRegHandler)
		authGroup.POST("/auth", h.auth.AuthHandler)
		authGroup.POST("/passReset", h.auth.PassResetHandler)
		authGroup.GET("/passReset/confirm/:token", h.auth.ConfirmPassResetHandler)
		authGroup.POST("/refresh", h.auth.RefreshHandler)
		authGroup.POST("/logout", h.auth.LogoutHandler)

		authGroup.GET("/reg", h.services.RenderRegForm)
		authGroup.GET("/auth", h.services.RenderAuthForm)
		authGroup.GET("/passReset", h.services.RenderPassResetForm)
		authGroup.GET("/changePass", h.services.RenderChangePassForm)
		authGroup.GET("/changeSkin", h.services.RenderChangeSkinForm)
		authGroup.GET("/deleteSkin", h.services.RenderDeleteSkinForm)
		authGroup.GET("/changeCloak", h.services.RenderChangeCloakForm)
		authGroup.GET("/deleteCloak", h.services.RenderDeleteCloakForm)
	}

	cabinetGroup := router.Group("/cabinet").Use(h.auth.AuthMiddleware())
	{
		cabinetGroup.POST("/changePass", h.cabinet.ChangePassHandler)
		cabinetGroup.POST("/changeSkin", h.cabinet.ChangeSkinHandler)
		cabinetGroup.POST("/deleteSkin", h.cabinet.DeleteSkinHandler)
		cabinetGroup.POST("/changeCloak", h.cabinet.ChangeCloakHandler)
		cabinetGroup.POST("/deleteCloak", h.cabinet.DeleteCloakHandler)
		cabinetGroup.POST("/transferMoney", h.cabinet.TransferMoneyHandler)
	}

	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/user", controllers.CreateUser)
		apiGroup.GET("/user/name/:name", controllers.GetUser)
		apiGroup.PUT("/user/exp", controllers.UpdateUserExp)
		apiGroup.PUT("/user/rank", controllers.UpdateUserRank)
		apiGroup.PUT("/user/playtime", controllers.UpdateUserPlaytime)
		apiGroup.PUT("/user/lastSeen", controllers.UpdateUserLastSeen)
		apiGroup.POST("/ban", controllers.BanUser)
		apiGroup.POST("/unban", controllers.UnbanUser)
		apiGroup.POST("/mute", controllers.MuteUser)
		apiGroup.POST("/unmute", controllers.UnmuteUser)
	}

	return router
}
