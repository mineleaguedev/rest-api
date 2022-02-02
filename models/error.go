package models

type Error struct {
	Code    int64  `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}
