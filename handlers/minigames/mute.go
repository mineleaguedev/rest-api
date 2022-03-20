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

func (h *Handler) PlayerMuteHandler(c *gin.Context) {
	var input models.PlayerMuteRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerMuteValues)
		return
	}

	if _, err := h.db.Exec("UPDATE `mutes` SET `status` = false WHERE `username` = ?", input.Username); err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBUnmutingPlayer, err)
		return
	}

	var muteTo sql.NullTime
	if input.Minutes != nil {
		muteTo = sql.NullTime{
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

	if _, err := h.db.Exec("INSERT INTO `mutes` (`username`, `mute_to`, `reason`, `admin`) VALUES (?, ?, ?, ?)",
		input.Username, muteTo, reason, input.Admin); err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBMutingPlayer, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

func (h *Handler) PlayerUnmuteHandler(c *gin.Context) {
	var input models.PlayerUnmuteRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerUnmuteValues)
		return
	}

	res, err := h.db.Exec("UPDATE `mutes` SET `status` = false WHERE `username` = ?", input.Username)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBUnmutingPlayer, err)
		return
	}

	amount, err := res.RowsAffected()
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrDBGettingRowsAffected, err)
		return
	}

	if amount <= 0 {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrPlayerIsNotMuted)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
