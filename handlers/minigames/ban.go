package minigames

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"time"
)

func (h *Handler) PlayerBanHandler(c *gin.Context) {
	var input models.PlayerBanRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerBanValues)
		return
	}

	if _, err := h.db.Exec("UPDATE `bans` SET `status` = false WHERE `username` = ?", input.Username); err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBUnbanningPlayer, err)
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

	if _, err := h.db.Exec("INSERT INTO `bans` (`username`, `ban_to`, `reason`, `admin`) VALUES (?, ?, ?, ?)",
		input.Username, banTo, reason, input.Admin); err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBBanningPlayer, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

func (h *Handler) PlayerUnbanHandler(c *gin.Context) {
	var input models.PlayerUnbanRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerUnbanValues)
		return
	}

	res, err := h.db.Exec("UPDATE `bans` SET `status` = false WHERE `username` = ?", input.Username)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBUnbanningPlayer, err)
		return
	}

	amount, err := res.RowsAffected()
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrDBGettingRowsAffected, err)
		return
	}

	if amount <= 0 {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrPlayerIsNotBanned)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
