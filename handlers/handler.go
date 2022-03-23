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
		authGroup.POST("/passReset", h.auth.PassResetHandler)
		authGroup.GET("/passReset/confirm/:token", h.auth.PassResetConfirmHandler)
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

	cabinetGroup := router.Group("/").Use(h.auth.AuthMiddleware())
	{
		cabinetGroup.POST("/changePass", h.cabinet.PassChangeHandler)
		cabinetGroup.POST("/changeSkin", h.cabinet.SkinChangeHandler)
		cabinetGroup.POST("/deleteSkin", h.cabinet.SkinDeleteHandler)
		cabinetGroup.POST("/changeCloak", h.cabinet.CloakChangeHandler)
		cabinetGroup.POST("/deleteCloak", h.cabinet.CloakDeleteHandler)
		cabinetGroup.POST("/transferMoney", h.cabinet.MoneyTransferHandler)
	}

	minigamesGroup := router.Group("/").Use(h.admin.ServerAdminAuthMiddleware())
	{
		minigamesGroup.POST("/player", h.minigames.PlayerCreateHandler)
		minigamesGroup.GET("/player/name/:name", h.minigames.PlayerGetHandler)
		minigamesGroup.PUT("/player/exp", h.minigames.PlayerExpUpdateHandler)
		minigamesGroup.PUT("/player/rank", h.minigames.PlayerRankUpdateHandler)
		minigamesGroup.PUT("/player/coins", h.minigames.PlayerCoinsUpdateHandler)
		minigamesGroup.PUT("/player/playtime", h.minigames.PlayerPlaytimeUpdateHandler)
		minigamesGroup.PUT("/player/lastSeen", h.minigames.PlayerLastSeenUpdateHandler)
		minigamesGroup.POST("/ban", h.minigames.PlayerBanHandler)
		minigamesGroup.POST("/unban", h.minigames.PlayerUnbanHandler)
		minigamesGroup.POST("/mute", h.minigames.PlayerMuteHandler)
		minigamesGroup.POST("/unmute", h.minigames.PlayerUnmuteHandler)
	}

	mapsGroup := router.Group("/")
	{
		mapsGroup.GET("/map", h.maps.MapsGetHandler).Use(h.admin.ServerAdminAuthMiddleware())
		mapsGroup.GET("/map/:minigame", h.maps.MiniGameMapsGetHandler).Use(h.admin.ServerAdminAuthMiddleware())
		mapsGroup.GET("/map/:minigame/:format", h.maps.MiniGameFormatMapsGetHandler).Use(h.admin.ServerAdminAuthMiddleware())
		mapsGroup.GET("/map/:minigame/:format/:map", h.maps.MapVersionsGetHandler).Use(h.admin.ServerAdminAuthMiddleware())
		mapsGroup.GET("/map/:minigame/:format/:map/:version/world", h.maps.MapWorldGetHandler).Use(h.admin.ServerAdminAuthMiddleware())
		mapsGroup.GET("/map/:minigame/:format/:map/:version/config", h.maps.MapConfigGetHandler).Use(h.admin.ServerAdminAuthMiddleware())

		mapsGroup.POST("/map", h.maps.MapUploadHandler).Use(h.admin.AdminAuthMiddleware())
	}

	pluginsGroup := router.Group("/")
	{
		pluginsGroup.GET("/plugin", h.plugins.PluginsGetHandler).Use(h.admin.ServerAdminAuthMiddleware())
		pluginsGroup.GET("/plugin/:name", h.plugins.PluginVersionsGetHandler).Use(h.admin.ServerAdminAuthMiddleware())
		pluginsGroup.GET("/plugin/:name/:version/", h.plugins.PluginGetHandler).Use(h.admin.ServerAdminAuthMiddleware())

		pluginsGroup.POST("/plugin", h.plugins.PluginUploadHandler).Use(h.admin.AdminAuthMiddleware())
	}

	adminGroup := router.Group("/admin").Use(h.admin.AdminAuthMiddleware())
	{
		adminGroup.GET("/server", h.admin.ServersGetHandler)
		adminGroup.POST("/server", h.admin.ServerAddHandler)
	}

	return router
}
