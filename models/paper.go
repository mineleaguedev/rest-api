package models

type PaperResponse struct {
	Success  bool     `json:"success"`
	Versions []string `json:"versions"`
}
