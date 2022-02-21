package models

type MuteRequest struct {
	Username string  `json:"username" binding:"required"`
	Minutes  *int64  `json:"minutes"`
	Reason   *string `json:"reason"`
	Admin    string  `json:"admin" binding:"required"`
}

type UnmuteRequest struct {
	Username string `json:"username" binding:"required"`
	Admin    string `json:"admin" binding:"required"`
}