package models

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
