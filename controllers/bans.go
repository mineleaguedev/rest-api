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

	if _, err := MiniGamesDB.Exec("UPDATE `bans` SET `status` = false WHERE `username` = ?", input.Username); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error setting user ban status to false",
		})
		log.Printf("Error setting user ban status to false: %s\n", err.Error())
		return
	}

	var banTo sql.NullTime
	if input.Minutes != nil {
		banTo = sql.NullTime{
			Time:  time.Now().Add(time.Duration(*input.Minutes) * time.Minute),
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

	if _, err := MiniGamesDB.Exec("INSERT INTO `bans` (`username`, `ban_to`, `reason`, `admin`) VALUES (?, ?, ?, ?)",
		input.Username, banTo, reason, input.Admin); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error banning user",
		})
		log.Printf("Error banning user: %s\n", err.Error())
	} else {
		c.JSON(http.StatusOK, models.Response{
			Success: true,
		})
	}
}

func UnbanUser(c *gin.Context) {
	var input models.UnbanRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Missing one or more fields " + err.Error(),
		})
		return
	}

	if res, err := MiniGamesDB.Exec("UPDATE `bans` SET `status` = false WHERE `username` = ?", input.Username); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error unbanning user",
		})
		log.Printf("Error unbanning user: %s\n", err.Error())
	} else {
		if amount, err := res.RowsAffected(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Success: false,
				Message: "Error getting rows affected amount",
			})
			log.Printf("Error getting rows affected amount: %s\n", err.Error())
		} else {
			if amount <= 0 {
				c.JSON(http.StatusBadRequest, models.Error{
					Success: false,
					Message: "User is not banned",
				})
			} else {
				c.JSON(http.StatusOK, models.Response{
					Success: true,
				})
			}
		}
	}
}
