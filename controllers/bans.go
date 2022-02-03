package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
	"time"
)

func BanUser(c *gin.Context) {
	var input models.BanRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Missing one or more fields " + err.Error(),
		})
		return
	}

	if _, err := DB.Exec("UPDATE `bans` SET `status` = false WHERE `username` = ?", input.Username); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error setting user ban status to false",
		})
		log.Printf("Error setting user ban status to false: %s\n", err.Error())
		return
	}

	var banTo sql.NullTime
	if input.BanTo != nil {
		banTo = sql.NullTime{
			Time:  time.Unix(*input.BanTo, 0),
			Valid: true,
		}
	}

	var reason sql.NullString
	if input.Reason != nil {
		reason = sql.NullString{
			String: *input.Reason,
			Valid:  true,
		}
	}

	if _, err := DB.Exec("INSERT INTO `bans` (`username`, `ban_to`, `reason`, `admin`) VALUES (?, ?, ?, ?)",
		input.Username, banTo, reason, input.Admin); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error banning user",
		})
		log.Printf("Error banning user: %s\n", err.Error())
	} else {
		c.JSON(http.StatusOK, models.BanResponse{
			Success: true,
		})
	}
}
