package models

type BanRequest struct {
	Username string  `json:"username" binding:"required"`
	BanTo    *int64  `json:"banTo"`
	Reason   *string `json:"reason"`
	Admin    string  `json:"admin" binding:"required"`
}

type UnbanRequest struct {
	Username string `json:"username" binding:"required"`
	Admin    string `json:"admin" binding:"required"`
}
