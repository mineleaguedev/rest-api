package models

type MuteRequest struct {
	Username string  `json:"username" binding:"required"`
	MuteTo   *int64  `json:"muteTo"`
	Reason   *string `json:"reason"`
	Admin    string  `json:"admin" binding:"required"`
}

type MuteResponse struct {
	Success bool `json:"success"`
}
