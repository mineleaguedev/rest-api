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

func MuteUser(c *gin.Context) {
	var input models.MuteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Missing one or more fields " + err.Error(),
		})
		return
	}

	if _, err := DB.Exec("UPDATE `mutes` SET `status` = false WHERE `username` = ?", input.Username); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error setting user mute status to false",
		})
		log.Printf("Error setting user mute status to false: %s\n", err.Error())
		return
	}

	var muteTo sql.NullTime
	if input.MuteTo != nil {
		muteTo = sql.NullTime{
			Time:  time.Unix(*input.MuteTo, 0),
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

	if _, err := DB.Exec("INSERT INTO `mutes` (`username`, `mute_to`, `reason`, `admin`) VALUES (?, ?, ?, ?)",
		input.Username, muteTo, reason, input.Admin); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error muting user",
		})
		log.Printf("Error muting user: %s\n", err.Error())
	} else {
		c.JSON(http.StatusOK, models.Response{
			Success: true,
		})
	}
}

func UnmuteUser(c *gin.Context) {
	var input models.UnmuteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Missing one or more fields " + err.Error(),
		})
		return
	}

	if _, err := DB.Exec("UPDATE `mutes` SET `status` = false WHERE `username` = ?", input.Username); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error unmuting user",
		})
		log.Printf("Error unmuting user: %s\n", err.Error())
	} else {
		c.JSON(http.StatusOK, models.Response{
			Success: true,
		})
	}
}
