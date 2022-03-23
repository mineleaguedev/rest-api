package admin

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
)

type serverAddRequest struct {
	Ip string `json:"ip" binding:"required"`
}

type serversResponse struct {
	Success bool     `json:"success"`
	Servers []string `json:"servers"`
}

func (h *Handler) ServersGetHandler(c *gin.Context) {
	rows, err := h.generalDB.Query("SELECT INET_NTOA(`ip`) FROM `servers`")
	if err != nil {
		h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrDBGettingServers)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrDBGettingServers)
			return
		}
	}(rows)

	var servers []string
	for rows.Next() {
		var server string
		if err := rows.Scan(&server); err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, errors.ErrDBGettingServers)
			return
		}
		servers = append(servers, server)
	}

	c.JSON(http.StatusOK, serversResponse{
		Success: true,
		Servers: servers,
	})
}

func (h *Handler) ServerAddHandler(c *gin.Context) {
	var input serverAddRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingAddServerValues)
		return
	}

	if _, err := h.generalDB.Exec("INSERT INTO `servers` (`ip`) VALUES (INET_ATON(?))", input.Ip); err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrDBServerAlreadyExists)
		} else {
			h.services.HandleInternalErr(c, errors.ErrDBAddingServer, err)
		}
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
