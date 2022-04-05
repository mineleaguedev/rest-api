package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/handlers/admin"
	"github.com/mineleaguedev/rest-api/handlers/auth"
	"github.com/mineleaguedev/rest-api/handlers/cabinet"
	"github.com/mineleaguedev/rest-api/handlers/maps"
	"github.com/mineleaguedev/rest-api/handlers/minigames"
	"github.com/mineleaguedev/rest-api/handlers/plugins"
	"github.com/mineleaguedev/rest-api/handlers/velocity"
	"github.com/mineleaguedev/rest-api/models"
	"github.com/mineleaguedev/rest-api/services"
)

type Handler struct {
	cabinet   *cabinet.Handler
	auth      *auth.Handler
	admin     *admin.Handler
	minigames *minigames.Handler
	maps      *maps.Handler
	plugins   *plugins.Handler
	velocity  *velocity.Handler
	services  *services.Service
}

func NewHandler(services *services.Service, middleware models.JWTMiddleware, generalDB, minigamesDB *sql.DB) *Handler {
	return &Handler{
		cabinet:   cabinet.NewHandler(services, generalDB),
		auth:      auth.NewHandler(services, middleware, generalDB),
		admin:     admin.NewHandler(services, generalDB, minigamesDB),
		minigames: minigames.NewHandler(services, minigamesDB),
		maps:      maps.NewHandler(services),
		plugins:   plugins.NewHandler(services),
		velocity:  velocity.NewHandler(services),
		services:  services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	authGroup := router.Group("/")
	{
		authGroup.POST("/reg", h.auth.RegHandler)
		authGroup.GET("/reg/confirm/:token", h.auth.RegConfirmHandler)
		authGroup.POST("/auth", h.auth.AuthHandler)
		authGroup.POST("/pass", h.auth.PassResetHandler)
		authGroup.GET("/pass/confirm/:token", h.auth.PassResetConfirmHandler)
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

	cabinetGroupAuth := router.Group("/cabinet").Use(h.auth.AuthMiddleware())
	{
		cabinetGroupAuth.PUT("/pass", h.cabinet.PassChangeHandler)
		cabinetGroupAuth.POST("/skin", h.cabinet.SkinChangeHandler)
		cabinetGroupAuth.DELETE("/skin", h.cabinet.SkinDeleteHandler)
		cabinetGroupAuth.POST("/cloak", h.cabinet.CloakChangeHandler)
		cabinetGroupAuth.DELETE("/cloak", h.cabinet.CloakDeleteHandler)
		cabinetGroupAuth.POST("/money", h.cabinet.MoneyTransferHandler)
	}

	minigamesGroupServerAdminAuth := router.Group("/").Use(h.admin.ServerAdminAuthMiddleware())
	{
		minigamesGroupServerAdminAuth.POST("/player", h.minigames.PlayerCreateHandler)
		minigamesGroupServerAdminAuth.GET("/player/name/:name", h.minigames.PlayerGetHandler)
		minigamesGroupServerAdminAuth.PUT("/player/exp", h.minigames.PlayerExpUpdateHandler)
		minigamesGroupServerAdminAuth.PUT("/player/rank", h.minigames.PlayerRankUpdateHandler)
		minigamesGroupServerAdminAuth.PUT("/player/coins", h.minigames.PlayerCoinsUpdateHandler)
		minigamesGroupServerAdminAuth.PUT("/player/playtime", h.minigames.PlayerPlaytimeUpdateHandler)
		minigamesGroupServerAdminAuth.PUT("/player/lastSeen", h.minigames.PlayerLastSeenUpdateHandler)
		minigamesGroupServerAdminAuth.POST("/ban", h.minigames.PlayerBanHandler)
		minigamesGroupServerAdminAuth.POST("/unban", h.minigames.PlayerUnbanHandler)
		minigamesGroupServerAdminAuth.POST("/mute", h.minigames.PlayerMuteHandler)
		minigamesGroupServerAdminAuth.POST("/unmute", h.minigames.PlayerUnmuteHandler)
	}

	mapsGroupServerAdminAuth := router.Group("/map").Use(h.admin.ServerAdminAuthMiddleware())
	{
		mapsGroupServerAdminAuth.GET("/", h.maps.MapsGetHandler)
		mapsGroupServerAdminAuth.GET("/:minigame", h.maps.MiniGameMapsGetHandler)
		mapsGroupServerAdminAuth.GET("/:minigame/:format", h.maps.MiniGameFormatMapsGetHandler)
		mapsGroupServerAdminAuth.GET("/:minigame/:format/:map", h.maps.MapVersionsGetHandler)
		mapsGroupServerAdminAuth.GET("/:minigame/:format/:map/:version/world", h.maps.MapWorldGetHandler)
		mapsGroupServerAdminAuth.GET("/:minigame/:format/:map/:version/config", h.maps.MapConfigGetHandler)
	}

	pluginsGroupServerAdminAuth := router.Group("/plugin").Use(h.admin.ServerAdminAuthMiddleware())
	{
		pluginsGroupServerAdminAuth.GET("/", h.plugins.PluginsGetHandler)
		pluginsGroupServerAdminAuth.GET("/:name", h.plugins.PluginVersionsGetHandler)
		pluginsGroupServerAdminAuth.GET("/:name/:version/", h.plugins.PluginGetHandler)
	}

	velocityGroupServerAdminAuth := router.Group("/velocity").Use(h.admin.ServerAdminAuthMiddleware())
	{
		velocityGroupServerAdminAuth.GET("/", h.velocity.VelocityVersionsGetHandler)
		velocityGroupServerAdminAuth.GET("/:version", h.velocity.VelocityGetHandler)
	}

	adminGroupAdminAuth := router.Group("/admin").Use(h.admin.AdminAuthMiddleware())
	{
		adminGroupAdminAuth.GET("/server", h.admin.ServersGetHandler)
		adminGroupAdminAuth.POST("/server", h.admin.ServerAddHandler)
		adminGroupAdminAuth.DELETE("/server", h.admin.ServerDeleteHandler)
		adminGroupAdminAuth.POST("/map", h.admin.MapUploadHandler)
		adminGroupAdminAuth.POST("/plugin", h.admin.PluginUploadHandler)
		adminGroupAdminAuth.POST("/velocity", h.admin.VelocityUploadHandler)
		adminGroupAdminAuth.POST("/paper", h.admin.PaperUploadHandler)
	}

	return router
}
