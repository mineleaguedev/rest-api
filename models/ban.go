package models

type BanRequest struct {
	Username string  `json:"username" binding:"required"`
	BanTo    *int64  `json:"banTo"`
	Reason   *string `json:"reason"`
	Admin    string  `json:"admin" binding:"required"`
}

type BanResponse struct {
	Success bool `json:"success"`
}
