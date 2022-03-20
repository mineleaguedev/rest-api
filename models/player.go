package models

type PlayerCreateRequest struct {
	Username string `json:"username" binding:"required,max=20"`
}

type PlayerUpdateExpRequest struct {
	Username string `json:"username" binding:"required"`
	Exp      int64  `json:"exp" binding:"required"`
}

type PlayerUpdateRankRequest struct {
	Username string  `json:"username" binding:"required"`
	Rank     *string `json:"rank"`
	RankTo   *int64  `json:"rankTo"`
}

type PlayerUpdateCoinsRequest struct {
	Username string `json:"username" binding:"required"`
	Coins    int64  `json:"coins" binding:"required"`
}

type PlayerUpdatePlaytimeRequest struct {
	Username string `json:"username" binding:"required"`
	Playtime int64  `json:"playtime" binding:"required"`
}

type PlayerUpdateLastSeenRequest struct {
	Username string `json:"username" binding:"required"`
	LastSeen int64  `json:"lastSeen" binding:"required"`
}

type Player struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Exp      int64     `json:"exp"`
	Rank     *string   `json:"rank"`
	Coins    int64     `json:"coins"`
	Playtime int64     `json:"playtime"`
	LastSeen int64     `json:"lastSeen"`
	Ban      *BanInfo  `json:"ban"`
	Mute     *MuteInfo `json:"mute"`
}

type BanInfo struct {
	To     *int64  `json:"to"`
	Reason *string `json:"reason"`
	Admin  string  `json:"admin"`
}

type MuteInfo struct {
	To     *int64  `json:"to"`
	Reason *string `json:"reason"`
	Admin  string  `json:"admin"`
}

type PlayerResponse struct {
	Success bool    `json:"success"`
	Player  *Player `json:"player"`
}
