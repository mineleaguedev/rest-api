package models

type UserCreateInput struct {
	Username string `json:"username" binding:"required,max=20"`
}

type UserExpUpdateInput struct {
	Username string `json:"username" binding:"required"`
	Exp      int64  `json:"exp" binding:"required"`
}

type UserRankUpdateInput struct {
	Username string `json:"username" binding:"required"`
	Rank     string `json:"rank" binding:"required"`
	RankTo   *int64 `json:"rankTo"`
}

type UserPlaytimeUpdateInput struct {
	Username string `json:"username" binding:"required"`
	Playtime int64  `json:"playtime" binding:"required"`
}

type User struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Exp      int64   `json:"exp"`
	Rank     *string `json:"rank"`
	Playtime int64   `json:"playtime"`
	LastSeen int64   `json:"lastSeen"`
	Ban      *Ban    `json:"ban"`
	Mute     *Ban    `json:"mute"`
}

type Ban struct {
	To     *int64  `json:"to"`
	Reason *string `json:"reason"`
	Admin  string  `json:"admin"`
}

type UserResponse struct {
	Success bool  `json:"success"`
	User    *User `json:"user"`
}
