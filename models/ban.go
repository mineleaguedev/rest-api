package models

type Ban struct {
	To     int64  `json:"to"`
	Reason string `json:"reason"`
	Admin  string `json:"admin" binding:"required"`
}
