package models

type VelocityResponse struct {
	Success  bool     `json:"success"`
	Versions []string `json:"versions"`
}
