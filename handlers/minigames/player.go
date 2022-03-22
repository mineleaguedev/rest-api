package minigames

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"time"
)

type playerCreateRequest struct {
	Username string `json:"username" binding:"required,max=20"`
}

type playerUpdateExpRequest struct {
	Username string `json:"username" binding:"required"`
	Exp      int64  `json:"exp" binding:"required"`
}

type playerUpdateRankRequest struct {
	Username string  `json:"username" binding:"required"`
	Rank     *string `json:"rank"`
	RankTo   *int64  `json:"rankTo"`
}

type playerUpdateCoinsRequest struct {
	Username string `json:"username" binding:"required"`
	Coins    int64  `json:"coins" binding:"required"`
}

type playerUpdatePlaytimeRequest struct {
	Username string `json:"username" binding:"required"`
	Playtime int64  `json:"playtime" binding:"required"`
}

type playerUpdateLastSeenRequest struct {
	Username string `json:"username" binding:"required"`
	LastSeen int64  `json:"lastSeen" binding:"required"`
}

func (h *Handler) PlayerCreateHandler(c *gin.Context) {
	var input playerCreateRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerCreateValues)
		return
	}

	res, err := h.db.Exec("INSERT INTO `players` (`username`) VALUES (?)", input.Username)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrPlayerAlreadyExists)
		} else {
			h.services.HandleInternalErr(c, errors.ErrMiniGamesDBCreatingPlayer, err)
		}
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrDBGettingLastInsertId, err)
		return
	}

	c.JSON(http.StatusOK, models.PlayerResponse{
		Success: true,
		Player: &models.Player{
			ID:       id,
			Username: input.Username,
			LastSeen: time.Now().Unix(),
		},
	})
}

func (h *Handler) PlayerGetHandler(c *gin.Context) {
	username := c.Param("name")

	var player = models.Player{Username: username}

	var rank sql.NullString
	var lastSeen sql.NullTime
	if err := h.db.QueryRow("SELECT `id`, `exp`, `rank`, `coins`, `playtime`, `last_seen` FROM `players` WHERE `username` = ?",
		username).Scan(&player.ID, &player.Exp, &rank, &player.Coins, &player.Playtime, &lastSeen); err != nil {
		if err == sql.ErrNoRows {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrPlayerDoesNotExist)
		} else {
			h.services.HandleInternalErr(c, errors.ErrMiniGamesDBGettingPlayer, err)
		}
		return
	}
	if rank.Valid {
		player.Rank = &rank.String
	}
	player.LastSeen = lastSeen.Time.Unix()

	// ban info
	var ban models.BanInfo
	var banTo sql.NullTime
	if err := h.db.QueryRow("SELECT `ban_to`, `reason`, `admin` FROM `bans` WHERE `username` = ? AND `status` = true",
		username).Scan(&banTo, &ban.Reason, &ban.Admin); err != nil && err != sql.ErrNoRows {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBGettingPlayerBanInfo, err)
		return
	}
	if banTo.Valid {
		banToInt64 := banTo.Time.Unix()
		player.Ban.To = &banToInt64
		player.Ban = &ban
	}

	// mute info
	var mute models.MuteInfo
	var muteTo sql.NullTime
	if err := h.db.QueryRow("SELECT `mute_to`, `reason`, `admin` FROM `mutes` WHERE `username` = ? AND `status` = true",
		username).Scan(&muteTo, &mute.Reason, &mute.Admin); err != nil && err != sql.ErrNoRows {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBGettingPlayerMuteInfo, err)
		return
	}
	if muteTo.Valid {
		muteToInt64 := muteTo.Time.Unix()
		player.Mute.To = &muteToInt64
		player.Mute = &mute
	}

	c.JSON(http.StatusOK, models.PlayerResponse{
		Success: true,
		Player:  &player,
	})
}

func (h *Handler) PlayerExpUpdateHandler(c *gin.Context) {
	var input playerUpdateExpRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerUpdateExpValues)
		return
	}

	res, err := h.db.Exec("UPDATE `players` SET `exp` = `exp` + ? WHERE `username` = ?", input.Exp, input.Username)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBUpdatingPlayerExp, err)
		return
	}

	handlePlayerUpdateInfo(h, c, res)
}

func (h *Handler) PlayerRankUpdateHandler(c *gin.Context) {
	var input playerUpdateRankRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerUpdateRankValues)
		return
	}

	var rank sql.NullString
	if input.Rank != nil {
		rank = sql.NullString{
			String: *input.Rank,
			Valid:  true,
		}
	}

	var rankTo sql.NullTime
	if input.RankTo != nil {
		rankTo = sql.NullTime{
			Time:  time.Unix(*input.RankTo, 0),
			Valid: true,
		}
	}

	res, err := h.db.Exec("UPDATE `players` SET `rank` = ?, `rank_to` = ? WHERE `username` = ?", rank, rankTo, input.Username)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBUpdatingPlayerRank, err)
		return
	}

	handlePlayerUpdateInfo(h, c, res)
}

func (h *Handler) PlayerCoinsUpdateHandler(c *gin.Context) {
	var input playerUpdateCoinsRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerUpdateCoinsValues)
		return
	}

	res, err := h.db.Exec("UPDATE `players` SET `coins` = `coins` + ? WHERE `username` = ?", input.Coins, input.Username)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBUpdatingPlayerCoins, err)
		return
	}

	handlePlayerUpdateInfo(h, c, res)
}

func (h *Handler) PlayerPlaytimeUpdateHandler(c *gin.Context) {
	var input playerUpdatePlaytimeRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerUpdatePlaytimeValues)
		return
	}

	res, err := h.db.Exec("UPDATE `players` SET `playtime` = `playtime` + ? WHERE `username` = ?", input.Playtime, input.Username)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBUpdatingPlayerPlaytime, err)
		return
	}

	handlePlayerUpdateInfo(h, c, res)
}

func (h *Handler) PlayerLastSeenUpdateHandler(c *gin.Context) {
	var input playerUpdateLastSeenRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPlayerUpdateLastSeenValues)
		return
	}

	res, err := h.db.Exec("UPDATE `players` SET `last_seen` = ? WHERE `username` = ?", time.Unix(input.LastSeen, 0), input.Username)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrMiniGamesDBUpdatingPlayerLastSeen, err)
		return
	}

	handlePlayerUpdateInfo(h, c, res)
}

func handlePlayerUpdateInfo(h *Handler, c *gin.Context, res sql.Result) {
	amount, err := res.RowsAffected()
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrDBGettingRowsAffected, err)
		return
	}

	if amount <= 0 {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrPlayerDoesNotExist)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
