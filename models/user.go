package models

type UserInput struct {
	Username string `json:"username" binding:"required"`
}

type User struct {
	ID       int64  `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Exp      int64  `json:"exp" binding:"required"`
	Rank     string `json:"rank" binding:"required"`
	Playtime int64  `json:"playtime" binding:"required"`
	LastSeen int64  `json:"lastSeen" binding:"required"`
	Ban      *Ban   `json:"ban"`
	Mute     *Ban   `json:"mute"`
}
