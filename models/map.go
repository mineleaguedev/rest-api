package models

type MiniGames struct {
	Name   string   `json:"name"`
	Format []Format `json:"formats"`
}

type Format struct {
	Format string `json:"format"`
	Map    []Map  `json:"maps"`
}

type Map struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

type MapsResponse struct {
	Success   bool        `json:"success"`
	MiniGames []MiniGames `json:"minigames"`
}
