package admin

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"net"
	"net/http"
	"strings"
)

func getClientIP(httpServer *http.Request) net.IP {
	var userIP string
	if len(httpServer.Header.Get("CF-Connecting-IP")) > 1 {
		userIP = httpServer.Header.Get("CF-Connecting-IP")
	} else if len(httpServer.Header.Get("X-Forwarded-For")) > 1 {
		userIP = httpServer.Header.Get("X-Forwarded-For")
	} else if len(httpServer.Header.Get("X-Real-IP")) > 1 {
		userIP = httpServer.Header.Get("X-Real-IP")
	} else {
		userIP = httpServer.RemoteAddr
		if strings.Contains(userIP, ":") {
			return net.ParseIP(strings.Split(userIP, ":")[0])
		} else {
			return net.ParseIP(userIP)
		}
	}
	return net.ParseIP(userIP)
}

func (h *Handler) ServerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c.Request).String()

		var exists bool
		if err := h.generalDB.QueryRow("SELECT 1 FROM `servers` WHERE `ip` = INET_ATON(?)", ip).Scan(&exists); err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		if !exists {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (h *Handler) AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessDetails, err := h.services.ExtractTokenMetadata(c.Request)
		if err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		userId, err := h.services.GetAuthSession(accessDetails)
		if err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		var username string
		if err := h.generalDB.QueryRow("SELECT `username` FROM `users` WHERE `id` = ?", userId).Scan(&username); err != nil {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		var rank sql.NullString
		if err := h.minigamesDB.QueryRow("SELECT `rank` FROM `players` WHERE `username` = ?", username).Scan(&rank); err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		if !rank.Valid {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		if rank.String != "ADMIN" {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (h *Handler) ServerAdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// server auth
		ip := c.ClientIP()
		if ip == "::1" {
			ip = "127.0.0.1"
		}

		var exists bool
		if err := h.generalDB.QueryRow("SELECT 1 FROM `servers` WHERE `ip` = INET_ATON(?)", ip).Scan(&exists); err == nil {
			if exists {
				c.Next()
				return
			}
		}

		// admin auth
		accessDetails, err := h.services.ExtractTokenMetadata(c.Request)
		if err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		userId, err := h.services.GetAuthSession(accessDetails)
		if err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		var username string
		if err := h.generalDB.QueryRow("SELECT `username` FROM `users` WHERE `id` = ?", userId).Scan(&username); err != nil {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		var rank sql.NullString
		if err := h.minigamesDB.QueryRow("SELECT `rank` FROM `players` WHERE `username` = ?", username).Scan(&rank); err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		if !rank.Valid {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		if rank.String != "ADMIN" {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrFuckYouBitch)
			c.Abort()
			return
		}

		c.Next()
	}
}
